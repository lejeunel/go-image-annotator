package list

import (
	"context"

	u "github.com/lejeunel/go-image-annotator/entities/user"
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

type FakeRepo struct {
	ErrOnCount bool
	ErrOnList  bool
	Err        error
	Count_     int
}

func (r *FakeRepo) Count() (int64, error) {
	if r.ErrOnCount {
		return 0, r.Err
	}
	return int64(r.Count_), nil
}

func (r *FakeRepo) List(req Request) ([]*u.User, error) {
	if r.ErrOnList {
		return nil, r.Err

	}

	result := []*u.User{}
	for range req.PageSize {
		usr := u.NewUser("the-id")
		result = append(result, &usr)
	}
	return result, nil

}

type FailingAuth struct {
}

func (f FailingAuth) ListUsers(ctx context.Context) error {
	return e.ErrAuth
}
