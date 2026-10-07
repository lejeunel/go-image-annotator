package update

import (
	"testing"

	rl "github.com/lejeunel/go-image-annotator/entities/role"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.RoleRepo{}, WithAuth(&fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestUpdateNonExistingRoleShouldFail(t *testing.T) {
	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&fk.RoleRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name, NewName: "new-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateRoleWithNameAlreadyTakenShouldFail(t *testing.T) {
	p := &FakePresenter{}
	currentRole := rl.NewRole(rl.NewRoleId(), "current-role")
	existingRole := rl.NewRole(rl.NewRoleId(), "existing-role")
	itr := New(&fk.RoleRepo{ExistingRoles: []rl.Role{existingRole, currentRole}})
	itr.Execute(t.Context(), Request{Name: currentRole.Name, NewName: existingRole.Name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	role := rl.NewRole(rl.NewRoleId(), "a-role")
	itr := New(&fk.RoleRepo{
		ExistingRoles: []rl.Role{role},
		ErrOnUpdate:   e.ErrInternal,
	})
	itr.Execute(t.Context(),
		Request{Name: role.Name, NewName: role.Name}, p)
	assert.True(t, p.GotInternalErr)
}

func TestUpdateWithInvalidMethodShouldFail(t *testing.T) {
	p := &FakePresenter{}
	role := rl.NewRole(rl.NewRoleId(), "a-role")
	itr := New(&fk.RoleRepo{ExistingRoles: []rl.Role{role}},
		WithAuth(&fk.Auth{ExistingMethods: []string{"a-method"}}))
	itr.Execute(t.Context(),
		Request{
			Name: role.Name, NewName: role.Name,
			NewMethods: []string{"a-non-existing-method"},
		}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateRole(t *testing.T) {
	p := &FakePresenter{}
	currentRole := rl.NewRole(rl.NewRoleId(), "a-role")
	repo := &fk.RoleRepo{ExistingRoles: []rl.Role{currentRole}}
	itr := New(repo, WithAuth(&fk.Auth{ExistingMethods: []string{"a-method"}}))
	description := "updated-description"
	req := Request{
		Name:           currentRole.Name,
		NewName:        "updated-name",
		NewDescription: &description,
		NewMethods:     []string{"a-method"},
	}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.NewName, p.Got.Name)
	assert.Equal(t, req.NewDescription, p.Got.Description)
	assert.Equal(t, req.NewMethods, p.Got.Methods)
}
