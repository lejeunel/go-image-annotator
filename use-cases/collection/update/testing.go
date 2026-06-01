package update

import (
	"context"
	"slices"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Names []string
	Got   Model
	Err   error
}

func (r *FakeRepo) Update(m Model) error {
	if r.Err != nil {
		return r.Err
	}
	r.Got = m
	return nil
}
func (r *FakeRepo) Exists(n string) (bool, error) {
	if slices.Contains(r.Names, n) {
		return true, nil
	}
	return false, nil
}
func (r *FakeRepo) Group(n string) (*string, error) {
	group := "my-group"
	return &group, nil
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

func (f FailingAuth) UpdateCollection(ctx context.Context, g string) error {
	return e.ErrAuth
}
