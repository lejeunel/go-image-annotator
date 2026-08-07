package delete

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessDeleteCollection(Response) {
	p.GotSuccess = true
}

type TestingTransactor struct {
	Repos
}

func (m *TestingTransactor) RunInTx(
	fn func(Repos) error,
) error {
	return fn(m.Repos)
}
