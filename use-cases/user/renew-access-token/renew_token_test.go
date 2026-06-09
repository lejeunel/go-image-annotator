package renew_token

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

func TestCreateWithTokenHash(t *testing.T) {
	token := "new-token"
	hash := []byte("new-hash")
	repo := &FakeRepo{}
	itr := New(repo, &FakeTokenGenerator{Token: token, Hash_: hash})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.Equal(t, token, p.Got.PersonalAccessToken)
	assert.Equal(t, hash, repo.Got.HashPAT)
}
