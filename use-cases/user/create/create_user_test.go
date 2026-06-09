package create

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{}, &FakeTokenGenerator{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleDuplicateError(t *testing.T) {
	userId := "user@example.com"
	itr := New(&FakeRepo{Ids: []string{userId}},
		&FakeTokenGenerator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: userId}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateWithTokenHash(t *testing.T) {
	token := "the-token"
	hash := []byte("the-hash")
	repo := &FakeRepo{}
	itr := New(repo, &FakeTokenGenerator{Token: token, Hash_: hash})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.Equal(t, token, p.Got.PersonalAccessToken)
	assert.Equal(t, hash, repo.Got.HashPAT)
}
