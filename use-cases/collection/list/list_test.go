package list

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	fk "github.com/lejeunel/go-image-annotator/use-cases/fakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ErrOnCount: e.ErrInternal})
	itr.Execute(t.Context(), pa.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestInvalidPageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{})
	itr.Execute(t.Context(), pa.PaginationParams{Page: -1}, p)
	assert.Equal(t, p.GotValidationErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.CollectionRepo{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), pa.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestListCollection(t *testing.T) {
	count := 3
	pageSize := 2
	page := int64(1)

	repo := &fk.CollectionRepo{Count_: count}
	p := &FakePresenter{}
	itr := New(repo)
	req := pa.PaginationParams{PageSize: pageSize, Page: page}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, len(p.Got.Collections), pageSize, "page size")
	assert.Equal(t, int(p.Got.Pagination.TotalRecords), count, "total records")
	assert.Equal(t, int(p.Got.Pagination.TotalPages), 2, "total pages")
	assert.Equal(t, p.Got.Pagination.Page, page, "page")
	assert.Equal(t, p.Got.Pagination.PageSize, pageSize, "page size")
}
