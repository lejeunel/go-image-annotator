package list

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestListLabel(t *testing.T) {
	count := 3
	pageSize := 2
	page := 1
	repo := &FakeRepo{Count_: count}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	itr.Execute(Request{PageSize: pageSize, Page: int64(page)}, p)

	assert.Equal(t, len(p.Got.Labels), pageSize, "page size")
	assert.Equal(t, int(p.Got.Pagination.TotalRecords), count, "total records")
	assert.Equal(t, int(p.Got.Pagination.TotalPages), 2, "total pages")
	assert.Equal(t, int(p.Got.Pagination.Page), page, "page")
	assert.Equal(t, p.Got.Pagination.PageSize, pageSize, "page size")
}
