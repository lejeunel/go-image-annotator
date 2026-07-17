package reset_forgotten_password

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingTokenShouldFail(t *testing.T) {
	repo := &fk.UserRepo{Missing: true}
	itr := New(repo, &fk.Tokenizer{}, &fk.Validator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Token: "the-token"}, p)
	assert.True(t, p.GotNotFoundErr)
}

func TestFindStateFromHash(t *testing.T) {
	hash := []byte("the-hash")
	repo := &fk.UserRepo{ReturnPasswordState: &u.ForgotPasswordState{}}
	itr := New(repo, &fk.Tokenizer{ReturnHash: hash}, &fk.Validator{})
	itr.Execute(t.Context(), Request{Token: "the-token"}, &FakePresenter{})
	assert.Equal(t, hash, repo.GotHash)
}

func TestMissingHashShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{Missing: true}, &fk.Tokenizer{}, &fk.Validator{})
	itr.Execute(t.Context(), Request{}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}

func TestPasswordMismatchShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{}, &fk.Tokenizer{}, &fk.Validator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "2"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrPasswordMismatch)
}

func TestInvalidPasswordShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{}, &fk.Tokenizer{}, &fk.Validator{Invalid: true})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrInvalidPassword)
}

func TestExpiredStateShouldFail(t *testing.T) {
	p := &FakePresenter{}
	expiresAt := time.Now()
	state := u.ForgotPasswordState{Id: "user@mail.com", ExpiresAt: &expiresAt}
	itr := New(&fk.UserRepo{ReturnPasswordState: &state}, &fk.Tokenizer{}, &fk.Validator{},
		WithClock(clockwork.NewFakeClockAt(expiresAt.Add(time.Hour))))
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrExpiredToken)
}

func TestHandleErrorOnUpdatePassword(t *testing.T) {
	p := &FakePresenter{}
	state := u.ForgotPasswordState{Id: "user@mail.com"}
	itr := New(&fk.UserRepo{ReturnPasswordState: &state, ErrOnUpdatePassword: e.ErrInternal}, &fk.Tokenizer{}, &fk.Validator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestUpdatePassword(t *testing.T) {
	p := &FakePresenter{}
	hash := []byte("the-hash")
	id := "user@mail.com"
	state := u.ForgotPasswordState{Id: id}
	repo := &fk.UserRepo{ReturnPasswordState: &state}
	itr := New(repo, &fk.Tokenizer{ReturnHash: hash}, &fk.Validator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, hash, repo.GotHash)
	assert.Equal(t, id, repo.GotId)
}

func TestShouldDeleteTokenAfterSuccessfulUpdate(t *testing.T) {
	p := &FakePresenter{}
	hash := []byte("the-hash")
	id := "user@mail.com"
	state := u.ForgotPasswordState{Id: id}
	repo := &fk.UserRepo{ReturnPasswordState: &state}
	itr := New(repo, &fk.Tokenizer{ReturnHash: hash}, &fk.Validator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, true, repo.DeletedPreviousTokens)
}
