package ingester

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"sync"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/shared/errors"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	ast "github.com/lejeunel/go-image-annotator/modules/file-store"
	kv "github.com/lejeunel/go-image-annotator/modules/string-validator"
	vv "github.com/lejeunel/go-image-annotator/modules/value-validator"
)

type ImageSpecsDetector interface {
	Detect(io.Reader) (*im.Specs, io.Reader, error)
}

type ArtefactRepo interface {
	Store(string, io.Reader) error
	Delete(string) error
}

type Repos struct {
	ImageRepo
	LabelRepo
	CollectionRepo
	AnnotationRepo
	MetaDataRepo
}
type Transactor interface {
	RunInTx(fn func(Repos) error) error
}
type ImageIngester struct {
	Hasher hash.Hash
	Repos
	Transactor
	ArtefactRepo
	ImageSpecsDetector
	KeyValidator   kv.Validator
	ValueValidator vv.Validator
	clockwork.Clock
	mu sync.Mutex
}

type Option func(*ImageIngester)

func WithClock(c clockwork.Clock) Option {
	return func(i *ImageIngester) {
		i.Clock = c
	}
}

// WithKeyValidator overrides the validator applied to metadata keys.
func WithKeyValidator(v kv.Validator) Option {
	return func(i *ImageIngester) {
		i.KeyValidator = v
	}
}

// WithValueValidator overrides the validator applied to metadata values.
func WithValueValidator(v vv.Validator) Option {
	return func(i *ImageIngester) {
		i.ValueValidator = v
	}
}

func New(imr ImageRepo, clr CollectionRepo,
	lr LabelRepo, ar AnnotationRepo, mr MetaDataRepo, tra Transactor,
	fileStore ast.FileStore, hasher hash.Hash, specsDetector ImageSpecsDetector, opts ...Option,
) *ImageIngester {
	i := &ImageIngester{
		Repos: Repos{
			ImageRepo: imr, LabelRepo: lr, CollectionRepo: clr,
			AnnotationRepo: ar, MetaDataRepo: mr,
		},
		Transactor:   tra,
		ArtefactRepo: fileStore, Hasher: hasher,
		ImageSpecsDetector: specsDetector,
		KeyValidator:       kv.NewNameValidator(),
		ValueValidator:     vv.BaseTypeValidator{},
		Clock:              clockwork.NewRealClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *ImageIngester) Ingest(r Request) (*Response, error) {
	errCtx := "ingesting image"
	collection, err := i.findCollectionByName(r.Collection)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	imageId := im.NewImageId()
	image, err := i.buildImage(imageId, *collection, r.Labels, r.BoundingBoxes,
		r.Polygons, r.Meta)
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

func (i *ImageIngester) storeRawData(image im.Image, reader io.Reader) (*[]byte, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	i.mu.Lock()
	i.Hasher.Reset()
	if _, err := i.Hasher.Write(data); err != nil {
		return nil, err
	}
	hash := i.Hasher.Sum(nil)
	i.mu.Unlock()

	if err := i.ensureDuplicateImageDoesNotExists(hash); err != nil {
		return nil, err
	}

	if err := i.ArtefactRepo.Store(image.Filename(), bytes.NewReader(data)); err != nil {
		return nil, err
	}

	return &hash, nil
}

func (i *ImageIngester) ingestImage(
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

	if err := tx.ImageRepo.AddToCollection(image.Id, image.Collection.Name); err != nil {
		return fmt.Errorf("adding image to collection: %w", err)
	}

	for _, label := range image.Labels {
		if err := tx.AnnotationRepo.AddImageLabel(
			image.Id,
			image.Collection.Name,
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
			image.Collection.Name,
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
			image.Collection.Name,
			poly,
			&authorId,
			&now,
		); err != nil {
			return fmt.Errorf("adding polygon: %w", err)
		}
	}

	for _, meta := range image.Meta {
		if err := tx.MetaDataRepo.Add(
			image.Collection.Name,
			image.Id,
			meta.Key,
			meta.Value,
		); err != nil {
			return fmt.Errorf("adding metadata with key %v: %w", meta.Key, err)
		}
	}
	return nil
}

func (i *ImageIngester) buildImage(id im.ImageId, collection clc.Collection, labelNames []string,
	bboxes []a.BoundingBoxRequest, polygons []a.PolygonRequest, meta []m.MetaData,
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
	if err := i.appendMeta(&image, meta); err != nil {
		return nil, err
	}

	return &image, nil
}

func (i *ImageIngester) appendLabels(image *im.Image, labelNames []string) error {
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

func (i *ImageIngester) appendBoundingBoxes(image *im.Image, bboxes []a.BoundingBoxRequest) error {
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

func (i *ImageIngester) appendPolygons(image *im.Image, polygons []a.PolygonRequest) error {
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

// appendMeta validates metadata before anything is stored, so an invalid
// key or value is rejected without touching the file store.
func (i *ImageIngester) appendMeta(image *im.Image, meta []m.MetaData) error {
	baseErr := fmt.Errorf("appending metadata")
	for _, md := range meta {
		if err := i.KeyValidator.Validate(md.Key); err != nil {
			return fmt.Errorf("%w: validating key %v: %w", baseErr, md.Key, err)
		}
		if err := i.ValueValidator.Validate(md.Value); err != nil {
			return fmt.Errorf("%w: validating value of key %v: %w", baseErr, md.Key, err)
		}
	}
	image.Meta = meta
	return nil
}

func (i *ImageIngester) findCollectionByName(name string) (*clc.Collection, error) {
	collection, err := i.CollectionRepo.Find(name)
	baseErr := fmt.Errorf("finding collection with name %v", name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return collection, nil
}

func (i *ImageIngester) findLabelByName(name string) (*lbl.Label, error) {
	baseErr := fmt.Errorf("fetching label by name %v", name)
	label, err := i.LabelRepo.FindLabel(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return label, nil
}

func (i *ImageIngester) ensureDuplicateImageDoesNotExists(hash []byte) error {
	baseErr := fmt.Errorf(
		"ensuring that duplicate image does not exist using hex hash %v",
		hex.EncodeToString(hash),
	)
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
