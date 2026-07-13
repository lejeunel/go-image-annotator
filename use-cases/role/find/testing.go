package find

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err  error
	Role rl.Role
}

func (r *FakeRepo) Find(name string) (*rl.Role, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if name == r.Role.Name {
		return &r.Role, nil
	}
	return nil, e.ErrNotFound

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
