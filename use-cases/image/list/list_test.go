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
		&fk.ImageStore{}, 1)
	itr.Execute(
		Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}},
		p,
	)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestSanitizePaginationParams(t *testing.T) {
	p := &FakePresenter{}
	repo := fk.ImageRepo{}
	defaultPageSize := 1
	itr := New(&repo,
		&fk.ImageStore{}, defaultPageSize)
	itr.Execute(
		Request{
			Filtering:        im.Filtering{},
			PaginationParams: pa.PaginationParams{PageSize: 0, Page: 0},
		},
		p,
	)
	assert.Equal(t, int64(1), repo.GotPagination.Page)
	assert.Equal(t, defaultPageSize, repo.GotPagination.PageSize)
}

func TestQueryPaginationParams(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.ImageStore{}, 20)
	r := Request{
		Filtering:        im.Filtering{},
		PaginationParams: pa.PaginationParams{Page: 1, PageSize: 2},
	}
	itr.Execute(r, p)
	pg := repo.GotPagination
	assert.Equal(t, int(pg.Page), int(r.Page), "page")
	assert.Equal(t, pg.PageSize, r.PageSize, "page size")
}

func TestHandleErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{ErrOnList: e.ErrInternal},
		&fk.ImageStore{}, 1)
	itr.Execute(
		Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}},
		p,
	)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrOnImageStoreRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{}, &fk.ImageStore{ErrOnFind: e.ErrInternal}, 1)
	itr.Execute(
		Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}},
		p,
	)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{ErrOnCount: e.ErrInternal}, &fk.ImageStore{}, 1)
	itr.Execute(
		Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}},
		p,
	)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestPaginationMetaData(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{Count_: 12}
	itr := New(repo, &fk.ImageStore{}, 10)
	r := Request{
		Filtering:        im.Filtering{},
		PaginationParams: pa.PaginationParams{Page: 1, PageSize: 10},
	}
	itr.Execute(r, p)
	pg := p.Got.Pagination
	assert.Equal(t, pg.Page, r.Page)
	assert.Equal(t, pg.PageSize, r.PageSize)
	assert.Equal(t, 12, int(pg.TotalRecords))
	assert.Equal(t, 2, int(pg.TotalPages))
}

func TestQueryOrderingParams(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.ImageStore{}, 1)
	ord := im.OrderingArgs{{Field: "ingested_at", Order: im.DescOrder}}

	r := Request{
		Filtering:        im.Filtering{},
		PaginationParams: pa.PaginationParams{PageSize: 1},
		OrderingArgs:     ord,
	}
	itr.Execute(r, p)
	assert.Equal(t, ord, repo.GotOrdering)
}

func TestListImages(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.ImageStore{}, 1)
	r := Request{Filtering: im.Filtering{}, PaginationParams: pa.PaginationParams{PageSize: 1}}
	itr.Execute(r, p)
	assert.Equal(t, r.PageSize, len(p.Got.Images))
}
