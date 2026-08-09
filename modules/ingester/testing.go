package ingester

import (
	"github.com/jonboulle/clockwork"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
)

type TestingTransactor struct {
	Repos
}

func (m *TestingTransactor) RunInTx(
	fn func(Repos) error,
) error {
	return fn(m.Repos)
}

func NewTestingRepos() Repos {
	return Repos{
		ImageRepo:      &fk.ImageRepo{},
		CollectionRepo: &fk.CollectionRepo{},
		LabelRepo:      &fk.LabelRepo{},
		AnnotationRepo: &fk.AnnotationRepo{},
	}
}

func NewTestingIngester(repos Repos, opts ...Option) *Ingester {
	i := &Ingester{
		Hasher:             &fk.Hasher{},
		Repos:              repos,
		Transactor:         &TestingTransactor{repos},
		ArtefactRepo:       &fk.FileStore{},
		ImageSpecsDetector: &fk.SpecsDetector{Return: im.Specs{MIMEType: "image/jpeg"}},
		Clock:              clockwork.NewFakeClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
