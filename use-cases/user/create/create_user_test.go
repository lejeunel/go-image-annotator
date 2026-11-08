package create

import (
	"testing"

	grp "github.com/lejeunel/go-image-annotator/entities/group"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.GroupRepo{}, &fk.RoleRepo{}, &fk.Tokenizer{},
		&fk.Tokenizer{},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleDuplicateError(t *testing.T) {
	userId := "user@example.com"
	itr := New(&fk.UserRepo{ExistingIds: []string{userId}}, &fk.GroupRepo{}, &fk.RoleRepo{},
		&fk.Tokenizer{},
		&fk.Tokenizer{},
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: userId}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateWithTokenHash(t *testing.T) {
	hash := []byte("the-hash")
	repo := &fk.UserRepo{}
	itr := New(repo, &fk.GroupRepo{}, &fk.RoleRepo{}, &fk.Tokenizer{ReturnHash: hash},
		&fk.Tokenizer{},
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.Equal(t, hash, repo.Created.HashPAT)
}

func TestCreateWithRandomPassword(t *testing.T) {
	repo := &fk.UserRepo{}
	passwordHash := []byte("a-password-hash")
	itr := New(repo, &fk.GroupRepo{}, &fk.RoleRepo{},
		&fk.Tokenizer{},
		&fk.Tokenizer{ReturnHash: passwordHash},
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.Equal(t, passwordHash, repo.Created.HashPassword)
}

func TestCreateWithPassword(t *testing.T) {
	password := "a-password"
	hash := []byte("the-hash")
	repo := &fk.UserRepo{}
	pwGenerator := &fk.Tokenizer{ReturnHash: hash}
	itr := New(repo, &fk.GroupRepo{}, &fk.RoleRepo{},
		&fk.Tokenizer{},
		pwGenerator,
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Password: &password}, p)
	assert.Equal(t, password, pwGenerator.Hashed)
	assert.Equal(t, hash, repo.Created.HashPassword)
}

func TestCreateWithRolesAndGroups(t *testing.T) {
	roles := []string{"a-role"}
	groups := []string{"a-group"}
	repo := &fk.UserRepo{}
	itr := New(repo,
		&fk.GroupRepo{Existing: []grp.Group{grp.NewGroup(grp.NewGroupId(), "a-group")}},
		&fk.RoleRepo{ExistingNames: roles},
		&fk.Tokenizer{}, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Roles: &roles, Groups: &groups}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, roles, repo.Created.Roles)
	assert.Equal(t, groups, repo.Created.Groups)
}

func TestMissingGroupShouldFail(t *testing.T) {
	groups := []string{"missing-group"}
	repo := &fk.UserRepo{}
	itr := New(repo, &fk.GroupRepo{}, &fk.RoleRepo{}, &fk.Tokenizer{}, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Groups: &groups}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
	assert.Nil(t, repo.Created)
}

func TestMissingRoleShouldFail(t *testing.T) {
	roles := []string{"missing-role"}
	repo := &fk.UserRepo{}
	itr := New(repo, &fk.GroupRepo{}, &fk.RoleRepo{}, &fk.Tokenizer{}, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Roles: &roles}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
	assert.Nil(t, repo.Created)
}

func TestHandleErrorOnCheckingGroupExistence(t *testing.T) {
	groups := []string{"a-group"}
	repo := &fk.UserRepo{}
	itr := New(repo, &fk.GroupRepo{ErrOnExists: e.ErrInternal}, &fk.RoleRepo{},
		&fk.Tokenizer{}, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Groups: &groups}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
	assert.Nil(t, repo.Created)
}

func TestHandleErrorOnCheckingRoleExistence(t *testing.T) {
	roles := []string{"a-role"}
	repo := &fk.UserRepo{}
	itr := New(repo, &fk.GroupRepo{}, &fk.RoleRepo{ErrOnExists: e.ErrInternal},
		&fk.Tokenizer{}, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Roles: &roles}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
	assert.Nil(t, repo.Created)
}
