package ingester

import (
	"errors"
	"fmt"
	"io"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/shared/errors"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	ast "github.com/lejeunel/go-image-annotator/modules/file-store"
	"hash"
)

type IImageSpecsDetector interface {
	Detect(io.Reader) (*im.ImageSpecs, io.Reader, error)
}
type Interface interface {
	Ingest(Request) (*Response, error)
}

type Repos struct {
	ImageRepo
	LabelRepo
	CollectionRepo
	AnnotationRepo
}
type UnitOfWork interface {
	// RunInTx runs fn inside a single transaction. Every store
	// in the Repos value executes against that transaction.
	RunInTx(fn func(Repos) error) error
}
type Ingester struct {
	Hasher hash.Hash
	Repos
	UnitOfWork
	ArtefactRepo       ast.Interface
	ImageSpecsDetector IImageSpecsDetector
	clock              clockwork.Clock
}

type Option func(*Ingester)

func WithClock(c clockwork.Clock) Option {
	return func(i *Ingester) {
		i.clock = c
	}
}

func New(imr ImageRepo, clr CollectionRepo,
	lr LabelRepo, ar AnnotationRepo, uow UnitOfWork,
	fileStore ast.Interface, hasher hash.Hash, specsDetector IImageSpecsDetector, opts ...Option) *Ingester {
	i := &Ingester{
		Repos:        Repos{imr, lr, clr, ar},
		UnitOfWork:   uow,
		ArtefactRepo: fileStore, Hasher: hasher,
		ImageSpecsDetector: specsDetector,
		clock:              clockwork.NewRealClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i Ingester) Ingest(r Request) (*Response, error) {
	errCtx := "ingesting image"
	collection, err := i.findCollectionByName(r.Collection)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	imageId := im.NewImageId()
	image, err := i.buildImage(imageId, *collection, r.Labels, r.BoundingBoxes)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	specs, reader, err := i.ImageSpecsDetector.Detect(r.Reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)

	}

	hash, err := i.ingestRawData(imageId, reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	specs.IngestedAt = i.clock.Now()
	if err := i.UnitOfWork.RunInTx(func(tx Repos) error {
		if err := i.ingestImage(tx, r.UserId, image, *hash, *specs); err != nil {
			i.ArtefactRepo.Delete(image.Id)
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	return &Response{ImageId: image.Id, Collection: collection.Name}, nil

}
func (i *Ingester) ingestRawData(id im.ImageId, reader io.Reader) (*[]byte, error) {

	tee := io.TeeReader(reader, i.Hasher)
	hash := i.Hasher.Sum(nil)

	if err := i.ensureDuplicateImageDoesNotExists(hash); err != nil {
		return nil, err
	}

	if err := i.ArtefactRepo.Store(id, tee); err != nil {
		return nil, err
	}

	return &hash, nil

}

func (i *Ingester) buildImage(id im.ImageId, collection clc.Collection, labelNames []string,
	bboxes []a.BoundingBoxRequest) (*im.Image, error) {
	image := im.NewImage(id, collection)

	if err := i.appendLabels(&image, labelNames); err != nil {
		return nil, err
	}

	if err := i.appendBoundingBoxes(&image, bboxes); err != nil {
		return nil, err
	}

	return &image, nil

}

func (i Ingester) appendLabels(image *im.Image, labelNames []string) error {
	for _, labelName := range labelNames {
		label, err := i.findLabelByName(labelName)
		if err != nil {
			return err
		}
		if err := image.AddLabel(*label); err != nil {
			return err
		}
	}
	return nil

}
func (i Ingester) appendBoundingBoxes(image *im.Image, bboxes []a.BoundingBoxRequest) error {
	baseErr := fmt.Errorf("appending bounding boxes")
	for _, bbox := range bboxes {
		label, err := i.findLabelByName(bbox.Label)
		if err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
		box_ := a.NewBoundingBox(a.NewAnnotationId(), bbox.Xc, bbox.Yc, bbox.Width, bbox.Height, *label)
		if err := image.AddBoundingBox(box_); err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
	}
	return nil

}

func (i Ingester) ingestImage(tx Repos, authorId u.UserId, image *im.Image, hash []byte, specs im.ImageSpecs) error {
	now := i.clock.Now()

	if err := tx.ImageRepo.AddImage(image.Id, hash, specs); err != nil {
		return fmt.Errorf("adding image: %w", err)
	}

	if err := tx.ImageRepo.AddToCollection(image.Id, image.Collection.Id); err != nil {
		return fmt.Errorf("adding image to collection: %w", err)
	}

	for _, label := range image.Labels {
		if err := tx.AnnotationRepo.AddImageLabel(image.Id, image.Collection.Id, label, &authorId, &now); err != nil {
			return fmt.Errorf("adding image label to collection: %w", err)
		}
	}

	for _, box := range image.BoundingBoxes {
		if err := tx.AnnotationRepo.AddBoundingBox(image.Id, image.Collection.Id, box, &authorId, &now); err != nil {
			return fmt.Errorf("adding bounding box: %w", err)
		}
	}
	return nil

}

func (i Ingester) findCollectionByName(name string) (*clc.Collection, error) {
	collection, err := i.CollectionRepo.FindCollectionByName(name)
	baseErr := fmt.Errorf("finding collection with name %v", name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return collection, nil

}

func (i Ingester) findLabelByName(name string) (*lbl.Label, error) {
	baseErr := fmt.Errorf("fetching label by name %v", name)
	label, err := i.LabelRepo.FindLabel(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return label, nil

}

func (i Ingester) ensureDuplicateImageDoesNotExists(hash []byte) error {

	baseErr := fmt.Errorf("ensuring that duplicate image does not exist using hash")
	duplicateId, err := i.ImageRepo.FindImageIdByHash(hash)
	if duplicateId != nil {
		return fmt.Errorf("%w: found duplicate image with id %v: %w", baseErr, *duplicateId, e.ErrDuplicate)
	}

	if errors.Is(err, e.ErrNotFound) {
		return nil
	}
	return err
}
