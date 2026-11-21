package ingest

import (
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	ing "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

// TestingIngester records the request it receives.
type TestingIngester struct {
	Got *ing.Request
}

func (i *TestingIngester) Ingest(r ing.Request) (*ing.Response, error) {
	i.Got = &r
	return &ing.Response{Collection: r.Collection}, nil
}

func NewTestingIngester() *TestingIngester {
	return &TestingIngester{}
}

func NewTestingInteractor(repo CollectionRepo, opts ...Option) *Interactor {
	i := &Interactor{
		Ingester:       NewTestingIngester(),
		Auth:           auth.NewVoidAuth(),
		CollectionRepo: repo,
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

func (p *FakePresenter) SuccessIngestImage(r ing.Response) {
	p.Got = &r
	p.GotSuccess = true
}
