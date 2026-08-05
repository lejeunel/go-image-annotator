package bootstrap

import (
	"testing"

	r "github.com/lejeunel/go-image-annotator/entities/role"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func TestNoopIfAdminRoleExists(t *testing.T) {
	roleRepo := &fk.RoleRepo{ExistingNames: []string{"admin"}}
	itr := New(&fk.UserRepo{}, roleRepo, &fk.FileStore{}, &fk.Tokenizer{}, &fk.StringValidator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotSuccess)
	assert.True(t, p.Got.Skipped)
}

func TestShouldCreateDefaultRoles(t *testing.T) {
	roleRepo := &fk.RoleRepo{}
	itr := New(&fk.UserRepo{}, roleRepo, &fk.FileStore{}, &fk.Tokenizer{}, &fk.StringValidator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.NotNil(t, roleRepo.Created)
	assert.Equal(t, len(r.DefaultRoleNames), len(roleRepo.Created))
}
func TestInvalidPassword(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.RoleRepo{}, &fk.FileStore{},
		&fk.Tokenizer{}, &fk.StringValidator{Invalid: true})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotValidationErr)
}

func TestShouldCreateAdminUser(t *testing.T) {
	userRepo := &fk.UserRepo{}
	pwHash := []byte("the-hash")
	tokenizer := &fk.Tokenizer{ReturnHash: pwHash}
	store := &fk.FileStore{}
	itr := New(userRepo, &fk.RoleRepo{}, store, tokenizer, &fk.StringValidator{})
	p := &FakePresenter{}
	id := "admin@mail.com"
	pw := "the-admin-password"
	itr.Execute(t.Context(), Request{id, pw}, p)
	assert.NotNil(t, userRepo.Created)
	assert.True(t, userRepo.Created.IsAdmin())
	assert.Equal(t, id, userRepo.Created.Id)
	assert.Equal(t, pw, tokenizer.Hashed)
	assert.Equal(t, pwHash, userRepo.Created.HashPassword)
	assert.True(t, len(store.GotData) > 0)
}
