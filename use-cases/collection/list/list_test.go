package list

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}

	itr := New(&fk.CollectionRepo{
		Existing:   []clc.Collection{clc.NewCollection(clc.NewCollectionId(), "a-collection")},
		ErrOnCount: e.ErrInternal,
	}, 1, 10)
	itr.Execute(t.Context(), pa.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ErrOnList: e.ErrInternal}, 1, 10)
	itr.Execute(t.Context(), pa.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestListCollection(t *testing.T) {
	pageSize := 1
	page := int64(1)

	repo := &fk.CollectionRepo{
		Existing: []clc.Collection{clc.NewCollection(clc.NewCollectionId(), "a-collection")},
	}
	p := &FakePresenter{}
	itr := New(repo, 1, 10)
	req := pa.PaginationParams{PageSize: pageSize, Page: page}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, pageSize, len(p.Got.Collections))
	assert.Equal(t, 1, int(p.Got.Pagination.TotalRecords))
	assert.Equal(t, 1, int(p.Got.Pagination.TotalPages))
	assert.Equal(t, page, p.Got.Pagination.Page)
	assert.Equal(t, pageSize, p.Got.Pagination.PageSize)
}
