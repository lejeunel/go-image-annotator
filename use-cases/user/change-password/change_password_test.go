package change_password

import (
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestPasswordMismatchShouldFail(t *testing.T) {
	p := &FakePresenter{}
	user := u.NewUser("user@example.com")
	itr := New(&fk.UserRepo{Return: &user}, &fk.Tokenizer{}, &fk.Validator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "2"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrPasswordMismatch)
}

func TestInvalidPasswordShouldFail(t *testing.T) {
	p := &FakePresenter{}
	user := u.NewUser("user@example.com")
	itr := New(&fk.UserRepo{Return: &user}, &fk.Tokenizer{}, &fk.Validator{Invalid: true})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrInvalidPassword)
}

func TestHandleErrorOnUpdatePassword(t *testing.T) {
	p := &FakePresenter{}
	user := u.NewUser("user@example.com")
	itr := New(&fk.UserRepo{Return: &user, ErrOnUpdatePassword: e.ErrInternal},
		&fk.Tokenizer{}, &fk.Validator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}
func TestFailWhenCurrentPasswordIsWrong(t *testing.T) {
	p := &FakePresenter{}
	user := u.NewUser("user@mail.com")
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo, &fk.Tokenizer{FailVerify: true}, &fk.Validator{})
	current := "asdf"
	itr.Execute(t.Context(),
		Request{Id: user.Id, CurrentPassword: current, FirstPassword: "1", SecondPassword: "1"}, p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrInvalidPassword)
}

func TestChangePassword(t *testing.T) {
	p := &FakePresenter{}
	hash := []byte("the-hash")
	user := u.NewUser("user@mail.com")
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo, &fk.Tokenizer{ReturnHash: hash}, &fk.Validator{})
	itr.Execute(t.Context(), Request{Id: user.Id, FirstPassword: "1", SecondPassword: "1"}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, hash, repo.GotHash)
	assert.Equal(t, user.Id, repo.GotId)
}
