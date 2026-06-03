package ingest

import (
	"context"
	"errors"
	"fmt"
	"io"

	e "github.com/lejeunel/go-image-annotator/shared/errors"

	"log/slog"

	ast "github.com/lejeunel/go-image-annotator/app/file-store"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"hash"
)

type IImageSpecsDetector interface {
	Detect(io.Reader) (*im.ImageSpecs, io.Reader, error)
}

type Interactor struct {
	Hasher             hash.Hash
	ImageRepo          ImageRepo
	CollectionRepo     CollectionRepo
	AnnotationRepo     AnnotationRepo
	LabelRepo          LabelRepo
	ArtefactRepo       ast.Interface
	Logger             *slog.Logger
	ImageSpecsDetector IImageSpecsDetector
	auth               Auth
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func NewInteractor(imageRepo ImageRepo, collectionRepo CollectionRepo,
	labelRepo LabelRepo, annotationRepo AnnotationRepo,
	fileStore ast.Interface, hasher hash.Hash, specsDetector IImageSpecsDetector, opts ...Option) *Interactor {
	i := &Interactor{ImageRepo: imageRepo, CollectionRepo: collectionRepo,
		AnnotationRepo: annotationRepo, LabelRepo: labelRepo,
		ArtefactRepo: fileStore, Hasher: hasher,
		ImageSpecsDetector: specsDetector,
		Logger:             logging.NewNoOpLogger(),
		auth:               auth.PassThroughAuth{},
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	collection, err := i.findCollectionByName(r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}
	if err := i.auth.IngestImage(ctx, collection.Group); err != nil {
		i.handleError(err, out)
		return
	}

	imageId := im.NewImageId()
	image, err := i.buildImage(imageId, *collection, r.Labels, r.BoundingBoxes)
	if err != nil {
		i.handleError(err, out)
		return
	}

	specs, reader, err := i.ImageSpecsDetector.Detect(r.Reader)
	if err != nil {
		i.handleError(err, out)
		return

	}

	hash, err := i.ingestRawData(imageId, reader)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.ingestImage(image, *hash, *specs); err != nil {
		i.handleError(err, out)
		i.ImageRepo.Delete(image.Id)
		i.ArtefactRepo.Delete(image.Id)
		return
	}

	out.Success(NewIngestionResponse(image))

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "ingesting image"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.Logger.Error(errCtx, "error", err)
	out.Error(err)
}
func (i *Interactor) ingestRawData(id im.ImageId, reader io.Reader) (*[]byte, error) {

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

func (i *Interactor) buildImage(id im.ImageId, collection clc.Collection, labelNames []string,
	bboxes []BoundingBoxRequest) (*im.Image, error) {
	image := im.NewImage(id, collection)

	if err := i.appendLabels(&image, labelNames); err != nil {
		return nil, err
	}

	if err := i.appendBoundingBoxes(&image, bboxes); err != nil {
		return nil, err
	}

	return &image, nil

}

func (i Interactor) appendLabels(image *im.Image, labelNames []string) error {
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
func (i Interactor) appendBoundingBoxes(image *im.Image, bboxes []BoundingBoxRequest) error {
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

func (i Interactor) ingestImage(image *im.Image, hash []byte, specs im.ImageSpecs) error {

	if err := i.ImageRepo.AddImage(image.Id, hash, specs); err != nil {
		return fmt.Errorf("adding image: %w", err)
	}

	if err := i.ImageRepo.AddToCollection(image.Id, image.Collection.Id); err != nil {
		return fmt.Errorf("adding image to collection: %w", err)
	}

	for _, label := range image.Labels {
		if err := i.AnnotationRepo.AddImageLabel(image.Id, image.Collection.Id, label); err != nil {
			return fmt.Errorf("adding image label to collection: %w", err)
		}
	}

	for _, box := range image.BoundingBoxes {
		if err := i.AnnotationRepo.AddBoundingBox(image.Id, image.Collection.Id, box); err != nil {
			return fmt.Errorf("adding bounding box: %w", err)
		}
	}
	return nil

}

func (i Interactor) findCollectionByName(name string) (*clc.Collection, error) {
	collection, err := i.CollectionRepo.FindCollectionByName(name)
	baseErr := fmt.Errorf("finding collection with name %v", name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return collection, nil

}

func (i Interactor) findLabelByName(name string) (*lbl.Label, error) {
	baseErr := fmt.Errorf("fetching label by name %v", name)
	label, err := i.LabelRepo.FindLabel(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return label, nil

}

func (i Interactor) ensureDuplicateImageDoesNotExists(hash []byte) error {

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
