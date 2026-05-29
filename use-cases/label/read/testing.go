package read

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Label l.Label
	Err   error
}

func (r *FakeRepo) FindLabelByName(name string) (*l.Label, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if name == r.Label.Name {
		return &r.Label, nil
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
