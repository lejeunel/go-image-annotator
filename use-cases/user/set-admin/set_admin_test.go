package set_admin

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user@example.com", true, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnSet(t *testing.T) {
	itr := New(&fk.UserRepo{ErrOnSetAdmin: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user@example.com", true, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestSetAdmin(t *testing.T) {
	repo := &fk.UserRepo{}
	itr := New(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user@example.com", true, p)
	assert.True(t, p.GotSuccess)
	assert.True(t, p.Got.IsAdmin)
	assert.True(t, repo.GotSetAdmin)
}
