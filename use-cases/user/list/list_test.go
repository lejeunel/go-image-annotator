package list

import (
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserStore{},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), pag.PaginationParams{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleListingError(t *testing.T) {
	itr := New(&fk.UserStore{ErrOnList: e.ErrInternal},
		WithAuth(fk.Auth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), pag.PaginationParams{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestListUsers(t *testing.T) {
	user := u.NewUser("user@example.com")
	store := &fk.UserStore{ExistingUsers: []u.User{user}}
	p := &FakePresenter{}
	itr := New(store)
	itr.Execute(t.Context(), pag.PaginationParams{PageSize: 2, Page: int64(1)}, p)

	assert.Equal(t, 1, len(p.Got.Users))
	assert.Equal(t, 1, int(p.Got.Pagination.TotalRecords))
	assert.Equal(t, 1, int(p.Got.Pagination.Page))
}
