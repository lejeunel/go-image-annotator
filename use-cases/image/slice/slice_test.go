package slice

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
	b := q.NewFilterParserBuilder()
	b.AddField("collection", schema.Is[string]())
	fv := b.Build()
	ov := q.NewOrderParserBuilder().AddField("ingested_at").Build()
	st := &fk.ImageStore{}
	return New(st, fv, ov, 1, 10)
}

func TestReturnPaginationInfo(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	wantPagination := pa.Pagination{Page: 1, PageSize: 1, TotalRecords: 1}
	store := &fk.ImageStore{ReturnPagination: wantPagination}
	itr.ImageStore = store
	itr.Execute(
		Request{},
		p,
	)
	assert.NoError(t, p.GotErr)
	assert.Equal(t, wantPagination.PageSize, p.Got.Pagination.PageSize)
	assert.Equal(t, wantPagination.Page, p.Got.Pagination.Page)
}

func TestInvalidFilterQueryShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	fv := fk.FilterValidator{Err: e.ErrValidation}
	itr.FilterValidator = &fv
	query := "i-dont-know-what-to-type-here"
	r := Request{FilterStr: query}
	itr.Execute(r, p)
	assert.Equal(t, query, fv.Got)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestInvalidOrderingStrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	ov := fk.FilterValidator{Err: e.ErrValidation}
	itr.OrderingValidator = &ov
	query := "i-dont-know-what-to-type-here"
	r := Request{OrderStr: query}
	itr.Execute(r, p)
	assert.Equal(t, query, ov.Got)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestHandleErrOnSlice(t *testing.T) {
	p := &FakePresenter{}
	itr := SetupList()
	store := &fk.ImageStore{ErrOnSlice: e.ErrInternal}
	itr.ImageStore = store
	query := "collection:my-collection"
	itr.Execute(
		Request{FilterStr: query, PaginationParams: pa.PaginationParams{PageSize: 1}},
		p,
	)
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
