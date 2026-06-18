package delete

import (
	"context"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err          error
	ErrOnDelete  bool
	Missing      bool
	IsPopulated_ bool
	ReturnGroup  string
}

func (r *FakeRepo) GroupOfCollection(string) (*string, error) {
	return &r.ReturnGroup, nil
}
func (r *FakeRepo) Delete(string) error {

	if r.ErrOnDelete {
		return r.Err
	}
	return nil
}

func (r *FakeRepo) Exists(c string) (*bool, error) {
	var exist = true
	if r.Missing {
		exist = false
		return &exist, nil
	}
	return &exist, nil
}

func (r *FakeRepo) IsPopulated(c string) (*bool, error) {
	res := true
	if r.IsPopulated_ {
		return &res, nil
	}
	res = false
	return &res, nil
}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) Success() {
	p.GotSuccess = true
}

type FailingAuth struct {
}

func (f FailingAuth) DeleteGroup(ctx context.Context) error {
	return e.ErrAuth
}
