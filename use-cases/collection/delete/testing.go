package delete

import (
	"context"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(Response) {
	p.GotSuccess = true
}

type FailingAuth struct {
}

func (f FailingAuth) DeleteCollection(ctx context.Context, g string) error {
	return e.ErrAuthorization
}
