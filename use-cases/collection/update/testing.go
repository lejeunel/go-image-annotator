package update

import (
	"context"
	"slices"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeGroupRepo struct {
	Return *string
	Err    error
}

func (r *FakeGroupRepo) GroupOfCollection(string) (*string, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return r.Return, nil
}

type FakeCollectionRepo struct {
	Names []string
	Got   Model
	Err   error
}

func (r *FakeCollectionRepo) Update(m Model) error {
	if r.Err != nil {
		return r.Err
	}
	r.Got = m
	return nil
}
func (r *FakeCollectionRepo) Exists(n string) (bool, error) {
	if slices.Contains(r.Names, n) {
		return true, nil
	}
	return false, nil
}
func (r *FakeCollectionRepo) GroupOfCollection(n string) (*string, error) {
	group := "my-group"
	return &group, nil
}

type FakePresenter struct {
	Got        Response
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessUpdateCollection(r Response) {
	p.GotSuccess = true
	p.Got = r
}

type FailingAuth struct {
}

func (f FailingAuth) UpdateCollection(ctx context.Context, g string) error {
	return e.ErrAuthorization
}
