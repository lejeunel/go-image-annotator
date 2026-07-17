package list

import (
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ErrOnCount: e.ErrInternal})
	itr.Execute(t.Context(), pag.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestInvalidPageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{})
	itr.Execute(t.Context(), pag.PaginationParams{Page: -1}, p)
	assert.Equal(t, p.GotValidationErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), pag.PaginationParams{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestList(t *testing.T) {
	count := int64(3)
	pageSize := 2
	page := int64(1)

	repo := &fk.RoleRepo{Count_: count}
	p := &FakePresenter{}
	itr := New(repo)
	req := pag.PaginationParams{PageSize: pageSize, Page: page}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, len(p.Got.Roles), pageSize, "page size")
	assert.Equal(t, p.Got.Pagination.TotalRecords, count, "total records")
	assert.Equal(t, int(p.Got.Pagination.TotalPages), 2, "total pages")
	assert.Equal(t, p.Got.Pagination.Page, page, "page")
	assert.Equal(t, p.Got.Pagination.PageSize, pageSize, "page size")
}
