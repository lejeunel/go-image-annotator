package ingester

import (
	"github.com/jonboulle/clockwork"
	fk "github.com/lejeunel/go-image-annotator/fakes"
)

func NewTestingIngester(opts ...Option) *Ingester {
	i := &Ingester{
		ImageRepo:          &fk.ImageRepo{},
		CollectionRepo:     &fk.CollectionRepo{},
		LabelRepo:          &fk.LabelRepo{},
		AnnotationRepo:     &fk.AnnotationRepo{},
		ArtefactRepo:       &fk.FileStore{},
		Hasher:             &fk.Hasher{},
		ImageSpecsDetector: &fk.SpecsDetector{},
		clock:              clockwork.NewFakeClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
