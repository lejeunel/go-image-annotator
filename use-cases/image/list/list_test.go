package list

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
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

func TestPaginate(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	images := []im.Image{
		im.NewImage(im.NewImageId(), collection),
		im.NewImage(im.NewImageId(), collection),
	}
	store := fk.ImageStore{ReturnPaginated: images, ReturnCount: int64(2)}
	itr := New(&store, 10, 10)
	itr.Execute(
		Request{
			PaginationParams: pa.PaginationParams{PageSize: 2, Page: 1},
		},
		p,
	)
	assert.NoError(t, p.GotErr)
	assert.Equal(t, images, p.Got.Images)
	assert.Equal(t, int64(2), p.Got.Pagination.TotalRecords)
}
