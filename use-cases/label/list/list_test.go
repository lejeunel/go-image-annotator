package list

import (
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), pag.PaginationParams{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnCount: e.ErrInternal})
	itr.Execute(t.Context(), pag.PaginationParams{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestListLabel(t *testing.T) {
	count := 3
	pageSize := 2
	page := 1
	repo := &fk.LabelRepo{Count_: int64(count)}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), pag.PaginationParams{PageSize: pageSize, Page: int64(page)}, p)

	assert.Equal(t, len(p.Got.Labels), pageSize, "page size")
	assert.Equal(t, int(p.Got.Pagination.TotalRecords), count, "total records")
	assert.Equal(t, int(p.Got.Pagination.TotalPages), 2, "total pages")
	assert.Equal(t, int(p.Got.Pagination.Page), page, "page")
	assert.Equal(t, p.Got.Pagination.PageSize, pageSize, "page size")
}
