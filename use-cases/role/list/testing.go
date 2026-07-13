package list

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err        error
	ErrOnCount bool
	ErrOnList  bool
	Count_     int64
}

func (r *FakeRepo) Count() (*int64, error) {
	count := int64(0)
	if r.ErrOnCount {
		return &count, r.Err
	}
	return &r.Count_, nil
}

func (r *FakeRepo) List(req Request) ([]*rl.Role, error) {
	if r.ErrOnList {
		return nil, r.Err
	}

	result := []*rl.Role{}
	for range req.PageSize {
		c := rl.NewRole(rl.NewRoleId(), "a-role")
		result = append(result, &c)
	}
	return result, nil
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
