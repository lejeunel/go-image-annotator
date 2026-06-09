package read

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	Err        error
	Collection clc.Collection
}

func (r *FakeRepo) FindCollectionByName(name string) (*clc.Collection, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if name == r.Collection.Name {
		return &r.Collection, nil
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
