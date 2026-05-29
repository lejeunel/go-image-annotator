package list

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
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
	st.AssertEqual(t, "page size", len(p.Got.Collections), pageSize)
	st.AssertEqual(t, "total records", p.Got.Pagination.TotalRecords, count)
	st.AssertEqual(t, "total pages", int(p.Got.Pagination.TotalPages), 2)
	st.AssertEqual(t, "page", p.Got.Pagination.Page, page)
	st.AssertEqual(t, "page size", p.Got.Pagination.PageSize, pageSize)
}
