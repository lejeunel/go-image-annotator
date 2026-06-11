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
func TestNonExistingUserShouldFail(t *testing.T) {
	repo := &FakeRepo{Missing: true}
	itr := New(repo, &FakeTokenGenerator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user"}, p)
	assert.True(t, p.GotNotFoundErr)
}

func TestCreateWithTokenHash(t *testing.T) {
	token := "new-token"
	hash := []byte("new-hash")
	repo := &FakeRepo{}
	itr := New(repo, &FakeTokenGenerator{Token: token, Hash_: hash})
	p := &FakePresenter{}
	id := "user"
	itr.Execute(t.Context(), Request{Id: id}, p)
	assert.Equal(t, token, p.Got.PersonalAccessToken)
	assert.Equal(t, hash, repo.GotHash)
	assert.Equal(t, id, repo.GotId)
}
