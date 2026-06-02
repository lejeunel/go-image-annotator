package list

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleNotFoundErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrNotFound},
		&st.FakeImageStore{})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnList: true, Err: e.ErrInternal},
		&st.FakeImageStore{})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnImageBuild(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrInternal})
	itr.Execute(Request{PageSize: 1}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal}, &st.FakeImageStore{})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestListImages(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := NewInteractor(repo, &st.FakeImageStore{})
	r := Request{Page: 1, PageSize: 2}
	itr.Execute(r, p)
	assert.Equal(t, r.PageSize, len(p.Got.Images))
}

func TestPaginationMetaData(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{Count_: 10}
	itr := NewInteractor(repo, &st.FakeImageStore{})
	r := Request{Page: 1, PageSize: 2}
	itr.Execute(r, p)
	pg := p.Got.Pagination
	assert.Equal(t, pg.Page, r.Page, "page")
	assert.Equal(t, pg.PageSize, r.PageSize, "page size")
	assert.Equal(t, int(pg.TotalRecords), 10, "total records")
	assert.Equal(t, int(pg.TotalPages), 5, "total pages")
}

func TestQueryCorrectPaginationWithFilters(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{Count_: 10}
	itr := NewInteractor(repo, &st.FakeImageStore{})
	r := Request{Page: 1, PageSize: 2}
	itr.Execute(r, p)
	f := repo.GotFilters
	assert.Equal(t, int(f.Page), int(r.Page), "page")
	assert.Equal(t, f.PageSize, r.PageSize, "page size")
}
