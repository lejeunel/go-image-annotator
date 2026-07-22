package user

import (
	"bytes"
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnSetTokenShouldFail(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	repo.Db.Close()
	err := repo.SetAccessTokenHash(userId, []byte(""))
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestSetAPIAccessTokenHash(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	CreateUser(repo, userId)
	hash := []byte("hello")
	err := repo.SetAccessTokenHash(userId, hash)
	assert.NoError(t, err)
	r, _ := repo.Find(userId)
	assert.True(t, bytes.Equal(r.HashPAT, hash))
}

func TestSetForgottenPasswordTokenHash(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	CreateUser(repo, userId)
	hash := []byte("hello")
	expiresAt := time.Now()
	err := repo.AddForgottenPasswordState(hash, userId, expiresAt)
	assert.NoError(t, err)
	r, err := repo.FindResetPasswordState(hash)
	assert.NoError(t, err)
	assert.Equal(t, userId, r.Id)
	assert.True(t, r.ExpiresAt.Equal(expiresAt))
}

func TestDeleteForgottenPasswordTokens(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	CreateUser(repo, userId)
	hash := []byte("hello")
	expiresAt := time.Now()
	repo.AddForgottenPasswordState(hash, userId, expiresAt)
	err := repo.DeleteForgottenPasswordTokens(userId)
	assert.NoError(t, err)
	_, err = repo.FindResetPasswordState(hash)
	assert.ErrorIs(t, err, e.ErrNotFound)
}
