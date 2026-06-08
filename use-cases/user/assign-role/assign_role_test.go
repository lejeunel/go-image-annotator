package assign_role

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := NewInteractor(&FakeRepo{UserMissing: true})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Role: "a-role"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAssignUserRoleAlreadyAssignedDoesNothing(t *testing.T) {
	roles := []string{"a-role"}
	user := usr.NewUser("user@example.com",
		usr.WithRoles(roles))
	repo := &FakeRepo{Return: &user}
	itr := NewInteractor(repo)
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
	repo := &FakeRepo{Return: &user}
	itr := NewInteractor(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Role: newGroup}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, updatedRoles, p.Got.Roles)
	assert.Equal(t, newGroup, *repo.GotNewRole)
}
