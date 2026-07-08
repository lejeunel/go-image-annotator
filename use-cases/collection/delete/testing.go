package delete

import (
	"context"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeCollectionRepo struct {
	Err          error
	ErrOnDelete  bool
	Missing      bool
	IsPopulated_ bool
}

type FakeGroupRepo struct {
	Return *string
	Err    error
}

func (r *FakeGroupRepo) GroupOfCollection(string) (*string, error) {
	return r.Return, nil
}
func (r *FakeCollectionRepo) Delete(string) error {

	if r.ErrOnDelete {
		return r.Err
	}
	return nil
}

func (r *FakeCollectionRepo) Exists(c string) (bool, error) {
	if r.Missing {
		return false, nil
	}
	return true, nil
}

func (r *FakeCollectionRepo) IsPopulated(c string) (*bool, error) {
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

func (f FailingAuth) DeleteCollection(ctx context.Context, g string) error {
	return e.ErrAuthorization
}
