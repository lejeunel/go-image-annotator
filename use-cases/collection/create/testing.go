package create

import (
	"slices"

	"context"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeGroupRepo struct {
	MissingGroup bool
	Return       string
}

func (r *FakeGroupRepo) Find(name string) (*grp.Group, error) {
	if r.MissingGroup {
		return nil, e.ErrNotFound
	}

	return &grp.Group{Name: r.Return}, nil
}

type FakeCollectionRepo struct {
	Err   error
	Names []string
	Got   clc.Collection
}

func (r *FakeCollectionRepo) Create(c clc.Collection) error {
	if r.Err != nil {
		return r.Err
	}

	r.Got = c
	return nil
}

func (r *FakeCollectionRepo) Exists(name string) (bool, error) {
	if slices.Contains(r.Names, name) {
		return true, nil
	}
	return false, nil
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

func (f FailingAuth) CreateCollection(ctx context.Context, g string) error {
	return e.ErrAuth
}
