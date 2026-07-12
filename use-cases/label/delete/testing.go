package delete

import (
	"context"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err         error
	IsUsed_     bool
	IsMissing   bool
	ErrOnDelete bool
	ErrOnIsUsed bool
	ErrOnExists bool
}

func (r *FakeRepo) Delete(string) error {
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (r *FakeRepo) IsUsed(n string) (*bool, error) {
	res := true
	if r.ErrOnIsUsed {
		res = false
		return &res, r.Err
	}
	if r.IsUsed_ {
		return &res, nil

	}
	res = false
	return &res, nil
}
func (r *FakeRepo) Exists(n string) (bool, error) {
	if r.ErrOnExists {
		return false, r.Err
	}
	if r.IsMissing {
		return false, nil
	}

	return true, nil

}

type FakePresenter struct {
	GotSuccess bool
	t.TestingErrPresenter
}

func (p *FakePresenter) SuccessDeleteLabel(string) {
	p.GotSuccess = true
}

type FailingAuth struct {
}

func (f FailingAuth) DeleteLabel(ctx context.Context) error {
	return e.ErrAuthorization
}
