package find

import (
	"context"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Return *u.User
	Err    error
}

func (r *FakeRepo) Find(id string) (*u.User, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	return r.Return, nil

}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(r Response) {
	p.GotSuccess = true
	p.Got = r
}

type FailingAuth struct {
}

func (f FailingAuth) FindUser(ctx context.Context) error {
	return e.ErrAuth
}
