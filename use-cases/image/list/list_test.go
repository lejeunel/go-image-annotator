package list

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestHandleNotFoundErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{ErrOnList: e.ErrNotFound},
		&fk.ImageStore{})
	itr.Execute(Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestPaginationParamsFallbackToDefault(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{},
		&fk.ImageStore{})
	itr.Execute(Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 0, Page: 0}}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{ErrOnList: e.ErrInternal},
		&fk.ImageStore{})
	itr.Execute(Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnImageBuild(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{}, &fk.ImageStore{Err: e.ErrInternal})
	itr.Execute(Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{ErrOnCount: e.ErrInternal}, &fk.ImageStore{})
	itr.Execute(Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestListImages(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.ImageStore{})
	r := Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}}
	itr.Execute(r, p)
	assert.Equal(t, r.PageSize, len(p.Got.Images))
}

func TestPaginationMetaData(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{Count_: 10}
	itr := New(repo, &fk.ImageStore{})
	r := Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{Page: 1, PageSize: 2}}
	itr.Execute(r, p)
	pg := p.Got.Pagination
	assert.Equal(t, pg.Page, r.Page, "page")
	assert.Equal(t, pg.PageSize, r.PageSize, "page size")
	assert.Equal(t, int(pg.TotalRecords), 10, "total records")
	assert.Equal(t, int(pg.TotalPages), 5, "total pages")
}

func TestQueryPaginationParams(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.ImageStore{})
	r := Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{Page: 1, PageSize: 2}}
	itr.Execute(r, p)
	pg := repo.GotPagination
	assert.Equal(t, int(pg.Page), int(r.Page), "page")
	assert.Equal(t, pg.PageSize, r.PageSize, "page size")
}

func TestQueryOrderingParams(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.ImageStore{})
	ord := im.Ordering{IngestTime: true}

	r := Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}, Ordering: ord}
	itr.Execute(r, p)
	assert.Equal(t, ord, repo.GotOrdering)
}
