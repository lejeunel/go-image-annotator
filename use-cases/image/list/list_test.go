package list

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	q "github.com/lejeunel/go-image-annotator/modules/query"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
	"go.tomakado.io/dumbql/schema"
)

func SetupList() Interactor {
	b := schema.NewSchemaBuilder()
	b.AddField("collection", schema.Is[string]())
	fv := q.NewFilterParser(b.Build())
	ov := q.NewOrderingConverter(q.WithOrderingField("ingested_at"))
	repo := &fk.ImageRepo{}
	st := &fk.ImageStore{}
	return New(repo, fv, ov, st, 1)
}

func TestSanitizePaginationParams(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	repo := fk.ImageRepo{}
	itr.Repo = &repo
	itr.Execute(
		Request{
			PaginationParams: pa.PaginationParams{PageSize: 0, Page: 0},
		},
		p,
	)
	assert.NoError(t, p.GotErr)
	assert.Equal(t, int64(1), repo.GotPagination.Page)
	assert.Equal(t, itr.DefaultPageSize, repo.GotPagination.PageSize)
}

func TestQueryPaginationParams(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	repo := fk.ImageRepo{}
	itr.Repo = &repo
	r := Request{
		PaginationParams: pa.PaginationParams{Page: 1, PageSize: 2},
	}
	itr.Execute(r, p)
	pg := repo.GotPagination
	assert.Equal(t, int(pg.Page), int(r.Page), "page")
	assert.Equal(t, pg.PageSize, r.PageSize, "page size")
}

func TestPaginationMetaData(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	itr.Repo = &fk.ImageRepo{Count_: 12}
	r := Request{
		PaginationParams: pa.PaginationParams{Page: 1, PageSize: 10},
	}
	itr.Execute(r, p)
	pg := p.Got.Pagination
	assert.Equal(t, pg.Page, r.Page)
	assert.Equal(t, pg.PageSize, r.PageSize)
	assert.Equal(t, 12, int(pg.TotalRecords))
	assert.Equal(t, 2, int(pg.TotalPages))
}

func TestInvalidFilterQueryShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	fv := fk.QueryStrValidator{Err: e.ErrValidation}
	itr.FilterQueryStrValidator = &fv
	query := "i-dont-know-what-to-type-here"
	r := Request{FilterStr: query}
	itr.Execute(r, p)
	assert.Equal(t, query, fv.Got)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestInvalidOrderingStrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	ov := fk.QueryStrValidator{Err: e.ErrValidation}
	itr.OrderingStrValidator = &ov
	query := "i-dont-know-what-to-type-here"
	r := Request{OrderStr: query}
	itr.Execute(r, p)
	assert.Equal(t, query, ov.Got)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestHandleErrOnSlice(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	repo := fk.ImageRepo{ErrOnSlice: e.ErrInternal}
	itr.Repo = &repo
	query := "collection:my-collection"
	itr.Execute(
		Request{FilterStr: query, PaginationParams: pa.PaginationParams{PageSize: 1}},
		p,
	)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
	assert.Equal(t, query, repo.GotFilters)
}

func TestHandleErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	itr.Repo = &fk.ImageRepo{ErrOnCount: e.ErrInternal}
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrOnImageStoreRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	itr.ImageStore = &fk.ImageStore{ErrOnFind: e.ErrInternal}
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestListImagesWithoutFilterAndOrder(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	r := Request{}
	itr.Execute(r, p)
	assert.True(t, p.GotSuccess)
}
