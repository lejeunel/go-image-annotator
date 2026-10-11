package list

import (
	"testing"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnList(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnList: e.ErrInternal})
	itr.Execute(t.Context(), pag.PaginationParams{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestHandleInternalErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnCount: e.ErrInternal})
	itr.Execute(t.Context(), pag.PaginationParams{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestListProfiles(t *testing.T) {
	pageSize := 2
	page := 1
	repo := &fk.ProfileRepo{Existing: []pr.Profile{pr.NewProfile(pr.NewProfileId(), "a-profile")}}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), pag.PaginationParams{PageSize: pageSize, Page: int64(page)}, p)

	assert.Equal(t, 1, len(p.Got.Profiles))
	assert.Equal(t, 1, int(p.Got.Pagination.TotalRecords))
	assert.Equal(t, 1, int(p.Got.Pagination.TotalPages))
	assert.Equal(t, page, int(p.Got.Pagination.Page))
	assert.Equal(t, pageSize, p.Got.Pagination.PageSize)
}
