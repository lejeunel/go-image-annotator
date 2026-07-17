package change_password

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.Tokenizer{},
		&fk.Validator{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
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

func TestHandleErrorOnUpdatePassword(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{ErrOnUpdatePassword: e.ErrInternal}, &fk.Tokenizer{}, &fk.Validator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestChangePassword(t *testing.T) {
	p := &FakePresenter{}
	hash := []byte("the-hash")
	id := "user@mail.com"
	repo := &fk.UserRepo{}
	itr := New(repo, &fk.Tokenizer{ReturnHash: hash}, &fk.Validator{})
	itr.Execute(t.Context(), Request{Id: id, FirstPassword: "1", SecondPassword: "1"}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, hash, repo.GotHash)
	assert.Equal(t, id, repo.GotId)
}
