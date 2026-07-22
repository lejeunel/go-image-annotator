package delete

import (
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user@example.com", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeletingMissingUserShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{})
	itr.Execute(t.Context(), "user@example.com", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}
func TestDeletingMyselfShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := "user@example.com"
	itr := New(&fk.UserRepo{ExistingIds: []string{id}})
	itr.Execute(u.AppendUserToContext(t.Context(), u.NewUser(id)), id, p)
	assert.True(t, p.GotForbiddenErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteUser(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{ExistingIds: []string{"user@example.com"}})
	itr.Execute(t.Context(), "user@example.com", p)
	assert.True(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{ExistingIds: []string{"user@example.com"},
		ErrOnDelete: e.ErrInternal})
	itr.Execute(t.Context(), "user@example.com", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
