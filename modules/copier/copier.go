package copier

import (
	"fmt"

	"github.com/jonboulle/clockwork"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	store "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type Interface interface {
	Copy(Request) (*Response, error)
}
type Copier struct {
	Store          store.Interface
	ImageRepo      ImageRepo
	CollectionRepo CollectionRepo
	AnnotationRepo AnnotationRepo
	clock          clockwork.Clock
}

type Option func(*Copier)

func WithClock(c clockwork.Clock) Option {
	return func(cp *Copier) {
		cp.clock = c
	}
}

func New(store store.Interface, collectionRepo CollectionRepo,
	annotationRepo AnnotationRepo, opts ...Option) *Copier {
	i := &Copier{Store: store, CollectionRepo: collectionRepo,
		AnnotationRepo: annotationRepo,
		clock:          clockwork.NewRealClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i Copier) Copy(r Request) (*Response, error) {
	errCtx := fmt.Errorf("copying image")

	image, err := i.Store.Find(im.BaseImage{ImageId: r.ImageId, Collection: r.SourceCollection})
	if err != nil {
		return nil, fmt.Errorf("%w: fetching source image: %w", errCtx, err)
	}

	dstCollection, err := i.CollectionRepo.FindCollectionByName(r.DestinationCollection)
	if err != nil {
		return nil, fmt.Errorf("%w: finding destination collection: %w", errCtx, err)
	}

	if r.Deep {
		for _, box := range image.BoundingBoxes {
			i.AnnotationRepo.AddBoundingBox(image.Id, dstCollection.Id,
				box, box.Author, box.Time)
		}

		for _, l := range image.Labels {
			i.AnnotationRepo.AddImageLabel(image.Id, dstCollection.Id,
				l, l.Author, l.Time)
		}

		for _, l := range image.Polygons {
			i.AnnotationRepo.AddPolygon(image.Id, dstCollection.Id,
				l, l.Author, l.Time)
		}
	}

	return nil, nil

}
