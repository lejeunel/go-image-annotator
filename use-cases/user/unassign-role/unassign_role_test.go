package unassign_role

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := New(&FakeRepo{Err: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestUnAssignUserWhoIsNotAssignedHasNoEffect(t *testing.T) {
	role := "a-role"
	user := usr.NewUser("user@example.com")
	repo := &FakeRepo{Return: &user}
	itr := New(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Role: role}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Roles))
	assert.Nil(t, repo.GotUnassignedRole)
}

func TestUnAssignUser(t *testing.T) {
	role := "a-role"
	user := usr.NewUser("user@example.com",
		usr.WithRoles([]string{role}))
	repo := &FakeRepo{Return: &user}
	itr := New(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Role: role}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Roles))
	assert.Equal(t, role, *repo.GotUnassignedRole)
}
