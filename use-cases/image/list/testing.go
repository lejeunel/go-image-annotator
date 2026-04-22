package list

import (
	ist "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	t "github.com/lejeunel/go-image-annotator-v2/shared/testing"
)

type FakeRepo struct {
	GotFilters ist.FilteringParams
	Err        error
	Count_     int64
	ErrOnList  bool
	ErrOnCount bool
}

func (r *FakeRepo) List(f ist.FilteringParams) (*[]im.BaseImage, error) {
	if r.ErrOnList {
		return nil, r.Err
	}

	r.GotFilters = f

	result := []im.BaseImage{}
	collectionName := "a-collection"
	for range f.PageSize {
		result = append(result, im.BaseImage{Collection: collectionName, ImageId: im.NewImageId()})
	}

	return &result, nil

}
func (r *FakeRepo) Count(f ist.CountingParams) (*int64, error) {
	if r.ErrOnCount {
		return nil, r.Err
	}
	return &r.Count_, nil

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
