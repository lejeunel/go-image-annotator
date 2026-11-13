package bootstrap

import (
	"testing"

	r "github.com/lejeunel/go-image-annotator/entities/role"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func Setup() (Interactor, FakePresenter) {
	return New(
			&fk.UserRepo{},
			&fk.RoleRepo{},
			&fk.Tokenizer{},
			&fk.StringValidator{},
			&fk.StringValidator{},
		),
		FakePresenter{}
}

func TestNoopIfAdminRoleExists(t *testing.T) {
	itr, presenter := Setup()
	roleRepo := &fk.RoleRepo{ExistingNames: []string{"admin"}}
	itr.RoleRepo = roleRepo
	itr.Execute(t.Context(), Request{}, &presenter)
	assert.True(t, presenter.GotSuccess)
	assert.True(t, presenter.Got.Skipped)
}

func TestShouldCreateDefaultRoles(t *testing.T) {
	roleRepo := &fk.RoleRepo{}
	itr, presenter := Setup()
	itr.RoleRepo = roleRepo
	itr.Execute(t.Context(), Request{}, &presenter)
	assert.NotNil(t, roleRepo.Created)
	assert.Equal(t, len(r.DefaultRoleNames), len(roleRepo.Created))
}

func TestInvalidPassword(t *testing.T) {
	itr, presenter := Setup()
	itr.PasswordValidator = &fk.StringValidator{Invalid: true}
	itr.Execute(t.Context(), Request{}, &presenter)
	assert.True(t, presenter.GotValidationErr)
}

func TestInvalidEmail(t *testing.T) {
	itr, presenter := Setup()
	itr.EmailValidator = &fk.StringValidator{Invalid: true}
	itr.Execute(t.Context(), Request{}, &presenter)
	assert.True(t, presenter.GotValidationErr)
}

func TestShouldCreateAdminUser(t *testing.T) {
	itr, presenter := Setup()
	userRepo := &fk.UserRepo{}
	pwHash := []byte("the-hash")
	tokenizer := &fk.Tokenizer{ReturnHash: pwHash}
	itr.UserRepo = userRepo
	itr.PasswordHasher = tokenizer
	id := "admin@mail.com"
	pw := "the-admin-password"
	itr.Execute(t.Context(), Request{id, pw}, &presenter)
	assert.NotNil(t, userRepo.Created)
	assert.True(t, userRepo.Created.IsAdmin())
	assert.Equal(t, id, userRepo.Created.Id)
	assert.Equal(t, pw, tokenizer.Hashed)
	assert.Equal(t, pwHash, userRepo.Created.HashPassword)
}
