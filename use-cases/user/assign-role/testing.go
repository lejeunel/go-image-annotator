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
	Err         error
	UserMissing bool
	Return      *usr.User
	GotNewRole  *string
}

func (r *FakeRepo) Find(id string) (*usr.User, error) {
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
func (r *FakeRepo) UserExists(id string) error {
	if r.UserMissing {
		return e.ErrNotFound
	}
	return nil
}

type FailingAuth struct {
}

func (f FailingAuth) AssignRoleToUser(ctx context.Context, id string, role string) error {
	return e.ErrAuth
}
