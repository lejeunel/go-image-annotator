package list

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{}, WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestInvalidPageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(t.Context(), Request{Page: -1}, p)
	assert.Equal(t, p.GotValidationErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{Page: 1, PageSize: 1}, p)
	assert.Equal(t, p.GotInternalErr, true)
	assert.Equal(t, p.GotSuccess, false)
}

func TestListCollection(t *testing.T) {
	count := int64(3)
	pageSize := 2
	page := int64(1)

	repo := &FakeRepo{Count_: count}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	req := Request{PageSize: pageSize, Page: page}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, len(p.Got.Collections), pageSize, "page size")
	assert.Equal(t, p.Got.Pagination.TotalRecords, count, "total records")
	assert.Equal(t, int(p.Got.Pagination.TotalPages), 2, "total pages")
	assert.Equal(t, p.Got.Pagination.Page, page, "page")
	assert.Equal(t, p.Got.Pagination.PageSize, pageSize, "page size")
}
