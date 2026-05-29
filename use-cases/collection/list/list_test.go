package list

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(Request{Page: 1, PageSize: 1}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error, but got none")
	}
}

func TestInvalidPageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(Request{Page: -1}, p)
	if !p.GotValidationErr || p.GotSuccess {
		t.Fatal("expected validation error")
	}
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrInternal})
	itr.Execute(Request{Page: 1, PageSize: 1}, p)
	if !p.GotInternalErr {
		t.Fatal("expected internal error, but got none")
	}
	if p.GotSuccess {
		t.Fatalf("expected to get no success")
	}
}

func TestListCollection(t *testing.T) {
	count := int64(3)
	pageSize := 2
	page := int64(1)

	repo := &FakeRepo{Count_: count}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	req := Request{PageSize: pageSize, Page: page}
	itr.Execute(req, p)
	assert.Equal(t, len(p.Got.Collections), pageSize, "page size")
	assert.Equal(t, p.Got.Pagination.TotalRecords, count, "total records")
	assert.Equal(t, int(p.Got.Pagination.TotalPages), 2, "total pages")
	assert.Equal(t, p.Got.Pagination.Page, page, "page")
	assert.Equal(t, p.Got.Pagination.PageSize, pageSize, "page size")
}
