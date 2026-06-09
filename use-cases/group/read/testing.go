package read

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err   error
	Group grp.Group
}

func (r *FakeRepo) Find(name string) (*grp.Group, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if name == r.Group.Name {
		return &r.Group, nil
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
