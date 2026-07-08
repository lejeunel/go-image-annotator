package create

import (
	"slices"

	"context"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err   error
	Names []string
	Got   *g.Group
}

func (r *FakeRepo) Create(g g.Group) error {
	if r.Err != nil {
		return r.Err
	}

	r.Got = &g
	return nil
}

func (r *FakeRepo) Exists(name string) (*bool, error) {
	var exist = true
	if slices.Contains(r.Names, name) {
		return &exist, nil
	}
	exist = false
	return &exist, nil
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

func (f FailingAuth) CreateGroup(ctx context.Context) error {
	return e.ErrAuthorization
}
