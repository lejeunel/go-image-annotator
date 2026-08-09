package ingester

import (
	"archive/zip"
	"errors"
	"fmt"
	"hash"
	"io"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/shared/errors"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	ast "github.com/lejeunel/go-image-annotator/modules/file-store"
)

type IImageSpecsDetector interface {
	Detect(io.Reader) (*im.Specs, io.Reader, error)
}

type Repos struct {
	ImageRepo
	LabelRepo
	CollectionRepo
	AnnotationRepo
}
type Transactor interface {
	RunInTx(fn func(Repos) error) error
}
type Ingester struct {
	Hasher hash.Hash
	Repos
	Transactor
	ArtefactRepo       ast.Interface
	ImageSpecsDetector IImageSpecsDetector
	clockwork.Clock
}

type Option func(*Ingester)

func WithClock(c clockwork.Clock) Option {
	return func(i *Ingester) {
		i.Clock = c
	}
}

func New(imr ImageRepo, clr CollectionRepo,
	lr LabelRepo, ar AnnotationRepo, tra Transactor,
	fileStore ast.Interface, hasher hash.Hash, specsDetector IImageSpecsDetector, opts ...Option,
) *Ingester {
	i := &Ingester{
		Repos:        Repos{imr, lr, clr, ar},
		Transactor:   tra,
		ArtefactRepo: fileStore, Hasher: hasher,
		ImageSpecsDetector: specsDetector,
		Clock:              clockwork.NewRealClock(),
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
	image, err := i.buildImage(imageId, *collection, r.Labels, r.BoundingBoxes,
		r.Polygons)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	specs, reader, err := i.ImageSpecsDetector.Detect(r.Reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	image.Specs = *specs
	hash, err := i.storeRawData(*image, reader)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	specs.IngestedAt = i.Clock.Now()
	if err := i.Transactor.RunInTx(func(tx Repos) error {
		if err := i.ingestImage(tx, r.UserId, image, *hash, *specs); err != nil {
			i.ArtefactRepo.Delete(image.Filename())
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	return &Response{ImageId: image.Id, Collection: collection.Name}, nil
}
func (i Ingester) IngestArchive(r BatchRequest) (*BatchResponse, error) {
	errCtx := fmt.Errorf("ingesting zip archive")

	numIngestedImages := int64(0)

	zr, err := zip.NewReader(r.ReaderAt, r.Size)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}

	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("%w: opening file: %w", errCtx, err)
		}

		_, err = i.Ingest(Request{UserId: r.UserId, Collection: r.Collection, Reader: reader})
		if err != nil {
			reader.Close()
			return nil, fmt.Errorf("%w: ingesting image: %w", errCtx, err)
		}

		if err := reader.Close(); err != nil {
			return nil, fmt.Errorf("%w: closing image reader: %w", errCtx, err)
		}
		numIngestedImages += 1
	}

	return &BatchResponse{NumIngestedImages: numIngestedImages}, nil
}
func (i *Ingester) storeRawData(image im.Image, reader io.Reader) (*[]byte, error) {
	tee := io.TeeReader(reader, i.Hasher)
	hash := i.Hasher.Sum(nil)

	if err := i.ensureDuplicateImageDoesNotExists(hash); err != nil {
		return nil, err
	}

	if err := i.ArtefactRepo.Store(image.Filename(), tee); err != nil {
		return nil, err
	}

	return &hash, nil
}

func (i Ingester) ingestImage(
	tx Repos,
	authorId u.UserId,
	image *im.Image,
	hash []byte,
	specs im.Specs,
) error {
	now := i.Clock.Now()

	if err := tx.ImageRepo.AddImage(image.Id, hash, specs); err != nil {
		return fmt.Errorf("adding image: %w", err)
	}

	if err := tx.ImageRepo.AddToCollection(image.Id, image.Collection.Id); err != nil {
		return fmt.Errorf("adding image to collection: %w", err)
	}

	for _, label := range image.Labels {
		if err := tx.AnnotationRepo.AddImageLabel(
			image.Id,
			image.Collection.Id,
			label,
			&authorId,
			&now,
		); err != nil {
			return fmt.Errorf("adding image label to collection: %w", err)
		}
	}

	for _, box := range image.BoundingBoxes {
		if err := tx.AnnotationRepo.AddBoundingBox(
			image.Id,
			image.Collection.Id,
			box,
			&authorId,
			&now,
		); err != nil {
			return fmt.Errorf("adding bounding box: %w", err)
		}
	}

	for _, poly := range image.Polygons {
		if err := tx.AnnotationRepo.AddPolygon(
			image.Id,
			image.Collection.Id,
			poly,
			&authorId,
			&now,
		); err != nil {
			return fmt.Errorf("adding polygon: %w", err)
		}
	}
	return nil
}

func (i *Ingester) buildImage(id im.ImageId, collection clc.Collection, labelNames []string,
	bboxes []a.BoundingBoxRequest, polygons []a.PolygonRequest,
) (*im.Image, error) {
	image := im.NewImage(id, collection)

	if err := i.appendLabels(&image, labelNames); err != nil {
		return nil, err
	}

	if err := i.appendBoundingBoxes(&image, bboxes); err != nil {
		return nil, err
	}
	if err := i.appendPolygons(&image, polygons); err != nil {
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
		box_ := a.NewBoundingBox(
			a.NewAnnotationId(),
			bbox.Xc,
			bbox.Yc,
			bbox.Width,
			bbox.Height,
			*label,
		)
		if err := image.AddBoundingBox(box_); err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
	}
	return nil
}

func (i Ingester) appendPolygons(image *im.Image, polygons []a.PolygonRequest) error {
	baseErr := fmt.Errorf("appending polygons")
	for _, p := range polygons {
		label, err := i.findLabelByName(p.Label)
		if err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
		polygon_ := a.NewPolygon(
			a.NewAnnotationId(),
			p.Points,
			*label,
		)
		if err := image.AddPolygon(polygon_); err != nil {
			return fmt.Errorf("%w: %w", baseErr, err)
		}
	}
	return nil
}

func (i Ingester) findCollectionByName(name string) (*clc.Collection, error) {
	collection, err := i.CollectionRepo.Find(name)
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
		return fmt.Errorf(
			"%w: found duplicate image with id %v: %w",
			baseErr,
			*duplicateId,
			e.ErrDuplicate,
		)
	}

	if errors.Is(err, e.ErrNotFound) {
		return nil
	}
	return err
}
