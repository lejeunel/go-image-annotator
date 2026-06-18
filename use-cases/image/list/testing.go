package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/shared/testing"
)

type FakeRepo struct {
	GotFilters  im.FilteringParams
	GotOrdering im.OrderingParams
	Err         error
	Count_      int64
	ErrOnList   bool
	ErrOnCount  bool
}

func (r *FakeRepo) List(f im.FilteringParams, o im.OrderingParams) ([]im.BaseImage, error) {
	if r.ErrOnList {
		return nil, r.Err
	}

	r.GotFilters = f
	r.GotOrdering = o

	result := []im.BaseImage{}
	collectionName := "a-collection"
	for range f.PageSize {
		result = append(result,
			im.BaseImage{
				Collection: collectionName,
				ImageId:    im.NewImageId()})
	}

	return result, nil

}
func (r *FakeRepo) Count(f im.CountingParams) (*int64, error) {
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
