package assign_role

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := New(&fk.UserRepo{Missing: true})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Role: "a-role"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := New(&fk.UserRepo{ErrOnFind: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAssignUserRoleAlreadyAssignedDoesNothing(t *testing.T) {
	roles := []string{"a-role"}
	user := usr.NewUser("user@example.com",
		usr.WithRoles(roles))
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Role: "a-role"}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, roles, p.Got.Roles)
	assert.Nil(t, repo.GotNewRole)
}

func TestAssignUser(t *testing.T) {
	user := usr.NewUser("user@example.com",
		usr.WithRoles([]string{"a-role"}))
	newGroup := "new-role"
	updatedRoles := []string{"a-role", newGroup}
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Role: newGroup}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, updatedRoles, p.Got.Roles)
	assert.Equal(t, newGroup, *repo.GotNewRole)
}
