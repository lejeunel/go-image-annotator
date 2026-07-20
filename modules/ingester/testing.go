package ingester

import (
	"github.com/jonboulle/clockwork"
	fk "github.com/lejeunel/go-image-annotator/fakes"
)

type memUoW struct {
	Repos
}

func (m *memUoW) RunInTx(
	fn func(Repos) error) error {
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
		UnitOfWork:         &memUoW{repos},
		ArtefactRepo:       &fk.FileStore{},
		ImageSpecsDetector: &fk.SpecsDetector{},
		clock:              clockwork.NewFakeClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
