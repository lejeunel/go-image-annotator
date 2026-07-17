package renew_token

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.Tokenizer{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}
func TestNonExistingUserShouldFail(t *testing.T) {
	repo := &fk.UserRepo{Missing: true}
	itr := New(repo, &fk.Tokenizer{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), "user", p)
	assert.True(t, p.GotNotFoundErr)
}

func TestCreateWithTokenHash(t *testing.T) {
	token := "new-token"
	hash := []byte("new-hash")
	repo := &fk.UserRepo{ExistingIds: []string{"user"}}
	itr := New(repo, &fk.Tokenizer{ReturnValue: token, ReturnHash: hash})
	p := &FakePresenter{}
	id := "user"
	itr.Execute(t.Context(), id, p)
	assert.Equal(t, token, p.Got.PersonalAccessToken)
	assert.Equal(t, hash, repo.GotHash)
	assert.Equal(t, id, repo.GotId)
}
