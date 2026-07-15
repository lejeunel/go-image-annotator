package list

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleNotFoundErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnList: true, Err: e.ErrNotFound},
		&st.FakeImageStore{})
	itr.Execute(Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestPaginationParamsFallbackToDefault(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{},
		&st.FakeImageStore{})
	itr.Execute(Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{PageSize: 0, Page: 0}}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnList: true, Err: e.ErrInternal},
		&st.FakeImageStore{})
	itr.Execute(Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnImageBuild(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrInternal})
	itr.Execute(Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal}, &st.FakeImageStore{})
	itr.Execute(Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestListImages(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := New(repo, &st.FakeImageStore{})
	r := Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{PageSize: 1}}
	itr.Execute(r, p)
	assert.Equal(t, r.PageSize, len(p.Got.Images))
}

func TestPaginationMetaData(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{Count_: 10}
	itr := New(repo, &st.FakeImageStore{})
	r := Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{Page: 1, PageSize: 2}}
	itr.Execute(r, p)
	pg := p.Got.Pagination
	assert.Equal(t, pg.Page, r.Page, "page")
	assert.Equal(t, pg.PageSize, r.PageSize, "page size")
	assert.Equal(t, int(pg.TotalRecords), 10, "total records")
	assert.Equal(t, int(pg.TotalPages), 5, "total pages")
}

func TestQueryPaginationParams(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := New(repo, &st.FakeImageStore{})
	r := Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{Page: 1, PageSize: 2}}
	itr.Execute(r, p)
	pg := repo.GotPagination
	assert.Equal(t, int(pg.Page), int(r.Page), "page")
	assert.Equal(t, pg.PageSize, r.PageSize, "page size")
}

func TestQueryOrderingParams(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := New(repo, &st.FakeImageStore{})
	ord := im.OrderingParams{IngestTime: true}

	r := Request{FilteringParams: im.FilteringParams{}, PaginationParams: im.PaginationParams{PageSize: 1}, OrderingParams: ord}
	itr.Execute(r, p)
	assert.Equal(t, ord, repo.GotOrdering)
}
