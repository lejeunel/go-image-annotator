package ingest

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ing "github.com/lejeunel/go-image-annotator/modules/ingester"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type TestingIngester struct{}

func (i TestingIngester) Ingest(ing.Request) (*ing.Response, error) {
	return nil, nil
}

func NewTestingIngester() TestingIngester {
	return TestingIngester{}
}

func NewTestingInteractor(repo CollectionRepo, opts ...Option) *Interactor {
	i := &Interactor{
		ingester: NewTestingIngester(),
		auth:     auth.NewVoidAuth(),
		repo:     repo,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

type FakePresenter struct {
	Got        *ing.Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(r ing.Response) {
	p.Got = &r
	p.GotSuccess = true
}
