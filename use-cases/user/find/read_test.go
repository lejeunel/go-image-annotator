package find

import (
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserStore{},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	store := fk.UserStore{ErrOnFind: e.ErrInternal}
	itr := New(&store)
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestFindUser(t *testing.T) {
	p := &FakePresenter{}
	user := u.NewUser("user@example.com")
	store := fk.UserStore{Return: &user}
	itr := New(&store)
	itr.Execute(t.Context(), "", p)
	assert.Equal(t, user, p.Got)
	assert.True(t, p.GotSuccess)
}
