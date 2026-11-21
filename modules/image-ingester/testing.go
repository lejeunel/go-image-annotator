package ingester

import (
	"github.com/jonboulle/clockwork"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	kv "github.com/lejeunel/go-image-annotator/modules/string-validator"
	vv "github.com/lejeunel/go-image-annotator/modules/value-validator"
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
		MetaDataRepo:   &fk.MetaDataRepo{},
	}
}

func NewTestingImageIngester(repos Repos, opts ...Option) *ImageIngester {
	i := &ImageIngester{
		Hasher:             &fk.Hasher{},
		Repos:              repos,
		Transactor:         &TestingTransactor{repos},
		ArtefactRepo:       &fk.FileStore{},
		ImageSpecsDetector: &fk.SpecsDetector{Return: im.Specs{MIMEType: "image/jpeg"}},
		KeyValidator:       kv.NewNameValidator(),
		ValueValidator:     vv.BaseTypeValidator{},
		Clock:              clockwork.NewFakeClock(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
