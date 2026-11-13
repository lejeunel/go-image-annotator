package create

import (
	"testing"

	grp "github.com/lejeunel/go-image-annotator/entities/group"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Setup() (Interactor, FakePresenter) {
	itr := New(&fk.UserRepo{}, &fk.GroupRepo{}, &fk.RoleRepo{}, &fk.Tokenizer{},
		&fk.Tokenizer{}, &fk.StringValidator{})
	p := FakePresenter{}
	return itr, p
}

func TestHandleAuthError(t *testing.T) {
	itr, presenter := Setup()
	itr.Auth = fk.Auth{ErrOnAuth: e.ErrAuthorization}
	itr.Execute(t.Context(), Request{}, &presenter)
	assert.True(t, presenter.GotAuthErr)
	assert.False(t, presenter.GotSuccess)
}

func TestHandleDuplicateError(t *testing.T) {
	userId := "user@example.com"
	itr, presenter := Setup()
	itr.UserRepo = &fk.UserRepo{ExistingIds: []string{userId}}
	itr.Execute(t.Context(), Request{Id: userId}, &presenter)
	assert.True(t, presenter.GotDuplicationErr)
	assert.False(t, presenter.GotSuccess)
}

func TestFailOnInvalidId(t *testing.T) {
	itr, presenter := Setup()
	itr.IdValidator = &fk.StringValidator{Invalid: true}
	itr.Execute(t.Context(), Request{Id: "an-id"}, &presenter)
	assert.True(t, presenter.GotValidationErr)
	assert.False(t, presenter.GotSuccess)
}

func TestCreateWithTokenHash(t *testing.T) {
	itr, presenter := Setup()
	hash := []byte("the-hash")
	userRepo := &fk.UserRepo{}
	itr.UserRepo = userRepo
	itr.APITokenGenerator = &fk.Tokenizer{ReturnHash: hash}
	itr.Execute(t.Context(), Request{Id: "user"}, &presenter)
	assert.Equal(t, hash, userRepo.Created.HashPAT)
}

func TestCreateWithRandomPassword(t *testing.T) {
	userRepo := &fk.UserRepo{}
	passwordHash := []byte("a-password-hash")
	itr, presenter := Setup()
	itr.PasswordGenerator = &fk.Tokenizer{ReturnHash: passwordHash}
	itr.UserRepo = userRepo
	itr.Execute(t.Context(), Request{Id: "user"}, &presenter)
	assert.Equal(t, passwordHash, userRepo.Created.HashPassword)
}

func TestCreateWithPassword(t *testing.T) {
	itr, presenter := Setup()
	password := "a-password"
	hash := []byte("the-hash")
	userRepo := &fk.UserRepo{}
	pwGenerator := &fk.Tokenizer{ReturnHash: hash}
	itr.UserRepo = userRepo
	itr.PasswordGenerator = pwGenerator
	itr.Execute(t.Context(), Request{Id: "user", Password: &password}, &presenter)
	assert.Equal(t, password, pwGenerator.Hashed)
	assert.Equal(t, hash, userRepo.Created.HashPassword)
}

func TestCreateWithRolesAndGroups(t *testing.T) {
	roles := []string{"a-role"}
	groups := []string{"a-group"}
	userRepo := &fk.UserRepo{}
	groupRepo := &fk.GroupRepo{Existing: []grp.Group{grp.NewGroup(grp.NewGroupId(), "a-group")}}
	roleRepo := &fk.RoleRepo{ExistingNames: roles}
	itr, presenter := Setup()
	itr.UserRepo = userRepo
	itr.GroupRepo = groupRepo
	itr.RoleRepo = roleRepo
	itr.Execute(t.Context(), Request{Id: "user", Roles: &roles, Groups: &groups}, &presenter)
	assert.True(t, presenter.GotSuccess)
	assert.Equal(t, roles, userRepo.Created.Roles)
	assert.Equal(t, groups, userRepo.Created.Groups)
}

func TestMissingGroupShouldFail(t *testing.T) {
	itr, presenter := Setup()
	groups := []string{"missing-group"}
	userRepo := &fk.UserRepo{}
	itr.UserRepo = userRepo
	itr.Execute(t.Context(), Request{Id: "user", Groups: &groups}, &presenter)
	assert.True(t, presenter.GotNotFoundErr)
	assert.False(t, presenter.GotSuccess)
	assert.Nil(t, userRepo.Created)
}

func TestMissingRoleShouldFail(t *testing.T) {
	itr, presenter := Setup()
	roles := []string{"missing-role"}
	userRepo := &fk.UserRepo{}
	itr.UserRepo = userRepo
	itr.Execute(t.Context(), Request{Id: "user", Roles: &roles}, &presenter)
	assert.True(t, presenter.GotNotFoundErr)
	assert.False(t, presenter.GotSuccess)
	assert.Nil(t, userRepo.Created)
}

func TestHandleErrorOnCheckingGroupExistence(t *testing.T) {
	itr, presenter := Setup()
	groups := []string{"a-group"}
	userRepo := &fk.UserRepo{}
	groupRepo := &fk.GroupRepo{ErrOnExists: e.ErrInternal}
	itr.UserRepo = userRepo
	itr.GroupRepo = groupRepo
	itr.Execute(t.Context(), Request{Id: "user", Groups: &groups}, &presenter)
	assert.True(t, presenter.GotInternalErr)
	assert.False(t, presenter.GotSuccess)
	assert.Nil(t, userRepo.Created)
}

func TestHandleErrorOnCheckingRoleExistence(t *testing.T) {
	itr, presenter := Setup()
	roles := []string{"a-role"}
	userRepo := &fk.UserRepo{}
	roleRepo := &fk.RoleRepo{ErrOnExists: e.ErrInternal}
	itr.UserRepo = userRepo
	itr.RoleRepo = roleRepo
	itr.Execute(t.Context(), Request{Id: "user", Roles: &roles}, &presenter)
	assert.True(t, presenter.GotInternalErr)
	assert.False(t, presenter.GotSuccess)
	assert.Nil(t, userRepo.Created)
}
