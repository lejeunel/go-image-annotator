package delete

import (
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err          error
	ErrOnDelete  bool
	Missing      bool
	IsPopulated_ bool
}

func (r *FakeRepo) Delete(string) error {

	if r.ErrOnDelete {
		return r.Err
	}
	return nil
}

func (r *FakeRepo) Exists(c string) (bool, error) {
	if r.Missing {
		return false, nil
	}
	return true, nil
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
