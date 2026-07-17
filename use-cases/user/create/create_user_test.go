package create

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.Tokenizer{},
		&fk.Tokenizer{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleDuplicateError(t *testing.T) {
	userId := "user@example.com"
	itr := New(&fk.UserRepo{ExistingIds: []string{userId}},
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
	itr := New(repo, &fk.Tokenizer{ReturnHash: hash},
		&fk.Tokenizer{},
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.Equal(t, hash, repo.Created.HashPAT)
}

func TestCreateWithRandomPassword(t *testing.T) {
	repo := &fk.UserRepo{}
	passwordHash := []byte("a-password-hash")
	itr := New(repo,
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
	itr := New(repo,
		&fk.Tokenizer{},
		pwGenerator,
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Password: &password}, p)
	assert.Equal(t, password, pwGenerator.GotToken)
	assert.Equal(t, hash, repo.Created.HashPassword)
}

func TestCreateAdmin(t *testing.T) {
	repo := &fk.UserRepo{}
	itr := New(repo, &fk.Tokenizer{}, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", IsAdmin: true}, p)
	assert.True(t, repo.Created.IsAdmin)
	assert.True(t, p.Got.IsAdmin)
}
