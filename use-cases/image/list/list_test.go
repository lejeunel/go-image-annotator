package list

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	stest "github.com/lejeunel/go-image-annotator/shared/testing"
)

func TestHandleNotFoundErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrNotFound},
		&st.FakeImageStore{})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrInternal},
		&st.FakeImageStore{})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected to get internal error")
	}
}

func TestHandleInternalErrOnImageBuild(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrInternal})
	itr.Execute(Request{PageSize: 1}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected to get internal error")
	}
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal}, &st.FakeImageStore{})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected to get internal error")
	}
}

func TestListImages(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := NewInteractor(repo, &st.FakeImageStore{})
	r := Request{Page: 1, PageSize: 2}
	itr.Execute(r, p)
	if !p.GotSuccess || (len(p.Got.Images) != r.PageSize) {
		t.Fatalf("expected to list images")
	}
}

func TestPaginationMetaData(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{Count_: 10}
	itr := NewInteractor(repo, &st.FakeImageStore{})
	r := Request{Page: 1, PageSize: 2}
	itr.Execute(r, p)
	pg := p.Got.Pagination
	stest.AssertEqual(t, "page", pg.Page, r.Page)
	stest.AssertEqual(t, "page size", pg.PageSize, r.PageSize)
	stest.AssertEqual(t, "total records", int(pg.TotalRecords), 10)
	stest.AssertEqual(t, "total pages", int(pg.TotalPages), 5)
}

func TestQueryCorrectPaginationWithFilters(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{Count_: 10}
	itr := NewInteractor(repo, &st.FakeImageStore{})
	r := Request{Page: 1, PageSize: 2}
	itr.Execute(r, p)
	f := repo.GotFilters
	stest.AssertEqual(t, "page", int(f.Page), int(r.Page))
	stest.AssertEqual(t, "page size", f.PageSize, r.PageSize)
}
