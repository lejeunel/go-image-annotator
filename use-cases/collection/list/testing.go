package list

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
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

func (r *FakeRepo) List(req Request) ([]*clc.Collection, error) {
	if r.ErrOnList {
		return nil, r.Err
	}

	result := []*clc.Collection{}
	for range req.PageSize {
		result = append(result, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
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
