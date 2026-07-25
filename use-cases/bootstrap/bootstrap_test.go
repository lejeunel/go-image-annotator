package bootstrap

import (
	"testing"

	r "github.com/lejeunel/go-image-annotator/entities/role"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func TestShouldSucceedIfAdminRoleExists(t *testing.T) {
	roleRepo := &fk.RoleRepo{ExistingNames: []string{"admin"}}
	itr := New(&fk.UserRepo{}, roleRepo, &fk.Tokenizer{}, &fk.Validator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotSuccess)
}

func TestShouldCreateDefaultRoles(t *testing.T) {
	roleRepo := &fk.RoleRepo{}
	itr := New(&fk.UserRepo{}, roleRepo, &fk.Tokenizer{}, &fk.Validator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.NotNil(t, roleRepo.Created)
	assert.Equal(t, len(r.DefaultRoleNames), len(roleRepo.Created))
}
func TestInvalidPassword(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.RoleRepo{},
		&fk.Tokenizer{}, &fk.Validator{Invalid: true})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotValidationErr)
}

func TestShouldCreateAdminUser(t *testing.T) {
	userRepo := &fk.UserRepo{}
	pwHash := []byte("the-hash")
	tokenizer := &fk.Tokenizer{ReturnHash: pwHash}
	itr := New(userRepo, &fk.RoleRepo{}, tokenizer, &fk.Validator{})
	p := &FakePresenter{}
	id := "admin@mail.com"
	pw := "the-admin-password"
	itr.Execute(t.Context(), Request{id, pw}, p)
	assert.NotNil(t, userRepo.Created)
	assert.True(t, userRepo.Created.IsAdmin())
	assert.Equal(t, id, userRepo.Created.Id)
	assert.Equal(t, pw, tokenizer.Hashed)
	assert.Equal(t, pwHash, userRepo.Created.HashPassword)
}
