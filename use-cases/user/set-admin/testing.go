package set_admin

import (
	"context"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success(r Response) {
	p.GotSuccess = true
	p.Got = r
}

type FakeRepo struct {
	Err      error
	Return   *u.User
	GotValue bool
}

func (r *FakeRepo) SetAdmin(id u.UserId, value bool) error {
	if r.Err != nil {
		return r.Err
	}
	r.GotValue = value
	return nil
}

type FailingAuth struct {
}

func (f FailingAuth) SetAdminRights(ctx context.Context) error {
	return e.ErrAuthorization
}
