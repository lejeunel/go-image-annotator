package unassign_group

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

type FakeGroupRepo struct {
	Missing bool
}

func (r *FakeGroupRepo) Exists(id string) (*bool, error) {
	if r.Missing {
		return nil, e.ErrNotFound
	}
	exist := true
	return &exist, nil
}

type FakeUserRepo struct {
	Err                error
	Missing            bool
	Return             *usr.User
	GotUnassignedGroup *string
}

func (r *FakeUserRepo) Find(id string) (*usr.User, error) {
	if r.Missing {
		return nil, e.ErrNotFound
	}
	return r.Return, nil
}
func (r *FakeUserRepo) UnAssignFromGroup(id string, group string) error {
	if r.Err != nil {
		return r.Err
	}
	r.GotUnassignedGroup = &group
	return nil
}

type FailingAuth struct {
}

func (f FailingAuth) UnAssignUserFromGroup(ctx context.Context, id string, group string) error {
	return e.ErrAuth
}
