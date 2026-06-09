package list

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestInvalidPageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{})
	itr.Execute(t.Context(), Request{Page: -1}, p)
	assert.Equal(t, p.GotValidationErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnList: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestList(t *testing.T) {
	count := int64(3)
	pageSize := 2
	page := int64(1)

	repo := &FakeRepo{Count_: count}
	p := &FakePresenter{}
	itr := New(repo)
	req := Request{PageSize: pageSize, Page: page}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, len(p.Got.Groups), pageSize, "page size")
	assert.Equal(t, p.Got.Pagination.TotalRecords, count, "total records")
	assert.Equal(t, int(p.Got.Pagination.TotalPages), 2, "total pages")
	assert.Equal(t, p.Got.Pagination.Page, page, "page")
	assert.Equal(t, p.Got.Pagination.PageSize, pageSize, "page size")
}
