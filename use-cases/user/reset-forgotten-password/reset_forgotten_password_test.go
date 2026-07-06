package reset_forgotten_password

import (
	"testing"
	"time"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingTokenShouldFail(t *testing.T) {
	repo := &FakeRepo{Missing: true}
	itr := New(repo, &FakeTokenHasher{}, &FakePasswordValidator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Token: "the-token"}, p)
	assert.True(t, p.GotNotFoundErr)
}

func TestFindStateFromHash(t *testing.T) {
	hash := []byte("the-hash")
	repo := &FakeRepo{Return: &u.ForgotPasswordState{}}
	itr := New(repo, &FakeTokenHasher{ReturnHash: hash}, &FakePasswordValidator{})
	itr.Execute(t.Context(), Request{Token: "the-token"}, &FakePresenter{})
	assert.Equal(t, hash, repo.GotHash)
}

func TestMissingHashShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Missing: true}, &FakeTokenHasher{}, &FakePasswordValidator{})
	itr.Execute(t.Context(), Request{}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}

func TestPasswordMismatchShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{}, &FakeTokenHasher{}, &FakePasswordValidator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "2"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrPasswordMismatch)
}

func TestInvalidPasswordShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{}, &FakeTokenHasher{}, &FakePasswordValidator{Invalid: true})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrInvalidPassword)
}

func TestHandleErrorOnUpdatePassword(t *testing.T) {
	p := &FakePresenter{}
	state := u.ForgotPasswordState{Id: "user@mail.com", ExpiresAt: time.Now()}
	itr := New(&FakeRepo{Return: &state, ErrOnUpdate: e.ErrInternal}, &FakeTokenHasher{}, &FakePasswordValidator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestUpdatePassword(t *testing.T) {
	p := &FakePresenter{}
	hash := []byte("the-hash")
	id := "user@mail.com"
	state := u.ForgotPasswordState{Id: id, ExpiresAt: time.Now()}
	repo := &FakeRepo{Return: &state}
	itr := New(repo, &FakeTokenHasher{ReturnHash: hash}, &FakePasswordValidator{})
	itr.Execute(t.Context(), Request{FirstPassword: "1", SecondPassword: "1"}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, hash, repo.GotHash)
	assert.Equal(t, id, repo.GotId)
}
