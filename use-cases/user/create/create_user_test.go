package create

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{}, &FakeAuthGenerator{},
		&FakeAuthGenerator{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleDuplicateError(t *testing.T) {
	userId := "user@example.com"
	itr := New(&FakeRepo{Ids: []string{userId}},
		&FakeAuthGenerator{},
		&FakeAuthGenerator{},
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: userId}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateWithTokenHash(t *testing.T) {
	hash := []byte("the-hash")
	repo := &FakeRepo{}
	itr := New(repo, &FakeAuthGenerator{Hash_: hash},
		&FakeAuthGenerator{},
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.Equal(t, hash, repo.Got.HashPAT)
}

func TestCreateWithRandomPassword(t *testing.T) {
	repo := &FakeRepo{}
	passwordHash := []byte("a-password-hash")
	itr := New(repo,
		&FakeAuthGenerator{},
		&FakeAuthGenerator{Hash_: passwordHash},
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.Equal(t, passwordHash, repo.Got.HashPassword)
}

func TestCreateWithPassword(t *testing.T) {
	password := "a-password"
	hash := []byte("the-hash")
	repo := &FakeRepo{}
	pwGenerator := &FakeAuthGenerator{Hash_: hash}
	itr := New(repo,
		&FakeAuthGenerator{},
		pwGenerator,
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", Password: &password}, p)
	assert.Equal(t, password, pwGenerator.GeneratedHashFromValue)
	assert.Equal(t, hash, repo.Got.HashPassword)
}

func TestCreateAdmin(t *testing.T) {
	repo := &FakeRepo{}
	itr := New(repo, &FakeAuthGenerator{}, &FakeAuthGenerator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user", IsAdmin: true}, p)
	assert.True(t, repo.Got.IsAdmin)
	assert.True(t, p.Got.IsAdmin)
}
