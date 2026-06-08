package assign_group

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
	Err          error
	UserMissing  bool
	GroupMissing bool
	Return       *usr.User
	GotNewGroup  *string
}

func (r *FakeRepo) Find(id string) (*usr.User, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return r.Return, nil
}
func (r *FakeRepo) AssignToGroup(id string, group string) error {
	if r.Err != nil {
		return r.Err
	}
	r.GotNewGroup = &group
	return nil
}
func (r *FakeRepo) UserExists(id string) error {
	if r.UserMissing {
		return e.ErrNotFound
	}
	return nil
}

func (r *FakeRepo) GroupExists(id string) error {
	if r.GroupMissing {
		return e.ErrNotFound
	}
	return nil
}

type FailingAuth struct {
}

func (f FailingAuth) AssignUserToGroup(ctx context.Context, id string, group string) error {
	return e.ErrAuth
}
