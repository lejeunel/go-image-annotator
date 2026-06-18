package assign_role

import (
	"context"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
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
	Err        error
	Missing    bool
	Return     *usr.User
	GotNewRole *string
}

func (r *FakeRepo) Find(id string) (*usr.User, error) {
	if r.Missing {
		return nil, e.ErrNotFound
	}
	if r.Err != nil {
		return nil, r.Err
	}
	return r.Return, nil
}
func (r *FakeRepo) AssignRole(id string, role string) error {
	if r.Err != nil {
		return r.Err
	}
	r.GotNewRole = &role
	return nil
}

type FailingAuth struct {
}

func (f FailingAuth) AssignRoleToUser(ctx context.Context) error {
	return e.ErrAuth
}
