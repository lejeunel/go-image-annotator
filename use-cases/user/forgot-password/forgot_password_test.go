package forgot_password

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingUserShouldFail(t *testing.T) {
	repo := &fk.UserRepo{Missing: true}
	itr := New(repo, 1, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user", p)
	assert.True(t, p.GotNotFoundErr)
}
func TestDeletePreviousTokens(t *testing.T) {
	repo := &fk.UserRepo{ExistingIds: []string{"user"}}
	itr := New(repo, 1, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user", p)
	assert.True(t, repo.DeletedPreviousTokens)
}

func TestRequestForgottenPasswordToken(t *testing.T) {
	token := "new-token"
	hash := []byte("new-hash")
	expiresInMinutes := 1
	repo := &fk.UserRepo{ExistingIds: []string{"user"}}

	now := time.Now()
	itr := New(repo,
		expiresInMinutes,
		&fk.Tokenizer{ReturnValue: token, ReturnHash: hash},
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
