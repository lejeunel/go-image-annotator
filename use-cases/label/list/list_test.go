package list

import (
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	st "github.com/lejeunel/go-image-annotator-v2/shared/testing"
	"testing"
)

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error, but got none")
	}
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error, but got none")
	}
}

func TestListLabel(t *testing.T) {
	count := 3
	pageSize := 2
	page := 1
	repo := &FakeRepo{Count_: count}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	itr.Execute(Request{PageSize: pageSize, Page: int64(page)}, p)

	st.AssertEqual(t, "page size", len(p.Got.Labels), pageSize)
	st.AssertEqual(t, "total records", int(p.Got.Pagination.TotalRecords), count)
	st.AssertEqual(t, "total pages", int(p.Got.Pagination.TotalPages), 2)
	st.AssertEqual(t, "page", int(p.Got.Pagination.Page), page)
	st.AssertEqual(t, "page size", p.Got.Pagination.PageSize, pageSize)
}
