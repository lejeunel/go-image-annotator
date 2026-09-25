package list

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"

	"github.com/stretchr/testify/assert"
)

func TestSanitizePaginationParams(t *testing.T) {
	p := &FakePresenter{}
	store := fk.ImageStore{}
	itr := New(&store, 10, 10)
	itr.Execute(
		Request{
			PaginationParams: pa.PaginationParams{PageSize: 0, Page: 0},
		},
		p,
	)
	assert.NoError(t, p.GotErr)
}

func TestHandleErrOnPaginate(t *testing.T) {
	p := &FakePresenter{}
	store := fk.ImageStore{ErrOnPaginateCollection: e.ErrInternal}
	itr := New(&store, 10, 10)
	itr.Execute(
		Request{},
		p,
	)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
	assert.False(t, p.GotSuccess)
}

func TestPaginate(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	images := []im.Image{
		im.NewImage(im.NewImageId(), collection),
		im.NewImage(im.NewImageId(), collection),
	}
	store := fk.ImageStore{
		ReturnPaginated:  images,
		ReturnPagination: pa.Pagination{Page: 1, PageSize: 1, TotalRecords: 1},
	}
	itr := New(&store, 10, 10)
	itr.Execute(
		Request{
			PaginationParams: pa.PaginationParams{PageSize: 2, Page: 1},
		},
		p,
	)
	assert.NoError(t, p.GotErr)
	assert.Equal(t, images, p.Got.Images)
	assert.Equal(t, int64(1), p.Got.Pagination.TotalRecords)
}
