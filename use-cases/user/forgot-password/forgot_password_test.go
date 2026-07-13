package forgot_password

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{},
		1,
		&FakeTokenGenerator{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}
func TestNonExistingUserShouldFail(t *testing.T) {
	repo := &FakeRepo{Missing: true}
	itr := New(repo, 1, &FakeTokenGenerator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user", p)
	assert.True(t, p.GotNotFoundErr)
}
func TestDeletePreviousTokens(t *testing.T) {
	repo := &FakeRepo{}
	itr := New(repo, 1, &FakeTokenGenerator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user", p)
	assert.True(t, repo.DeletedPreviousTokens)
}

func TestRequestForgottenPasswordToken(t *testing.T) {
	token := "new-token"
	hash := []byte("new-hash")
	expiresInMinutes := 1
	repo := &FakeRepo{}

	now := time.Now()
	itr := New(repo,
		expiresInMinutes,
		&FakeTokenGenerator{Token: token, Hash_: hash},
		WithClock(clockwork.NewFakeClockAt(now)))
	p := &FakePresenter{}
	email := "user"
	itr.Execute(t.Context(), email, p)
	assert.Equal(t, token, p.Got.PasswordResetToken)
	assert.Equal(t, email, p.Got.Email)
	assert.Equal(t, email, p.Got.Id)
	assert.Equal(t, hash, repo.GotHash)
	assert.Equal(t, email, repo.GotId)
	assert.Equal(t, now.Add(time.Minute*time.Duration(expiresInMinutes)), repo.GotExpiresAt)
}
