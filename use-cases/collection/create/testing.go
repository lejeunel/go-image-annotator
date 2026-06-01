package create

import (
	"slices"

	"context"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err   error
	Names []string
	Got   clc.Collection
}

func (r *FakeRepo) Create(c clc.Collection) error {
	if r.Err != nil {
		return r.Err
	}

	r.Got = c
	return nil
}

func (r *FakeRepo) Exists(name string) (bool, error) {
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
