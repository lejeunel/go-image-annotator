package update_role

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.RoleRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := New(&fk.UserRepo{Missing: true}, &fk.RoleRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Roles: []string{"a-role"}}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := New(&fk.UserRepo{ErrOnFind: e.ErrInternal}, &fk.RoleRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAssignUserRoleAlreadyAssignedDoesNothing(t *testing.T) {
	roles := []string{"a-role"}
	user := usr.NewUser("user@example.com",
		usr.WithRoles(roles))
	usrRepo := &fk.UserRepo{Return: &user}
	roleRepo := &fk.RoleRepo{ExistingNames: []string{"a-role"}}
	itr := New(usrRepo, roleRepo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Roles: []string{"a-role"}}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, roles, p.Got.Roles)
	assert.Nil(t, usrRepo.GotNewRole)
}

func TestAssignRoles(t *testing.T) {
	user := usr.NewUser("user@example.com",
		usr.WithRoles([]string{"a-role"}))
	newRole := "a-new-role"
	updatedRoles := []string{newRole}
	usrRepo := &fk.UserRepo{Return: &user}
	roleRepo := &fk.RoleRepo{ExistingNames: []string{"a-role", "a-new-role"}}
	itr := New(usrRepo, roleRepo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Roles: []string{newRole}}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, updatedRoles, usrRepo.SetRoles_)
	assert.Equal(t, user.Id, usrRepo.SetRolesToUser)
}
