package user

import (
	"bytes"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnSetTokenShouldFail(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	repo.Db.Close()
	err := repo.SetAccessTokenHash(userId, []byte(""))
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestSetTokenHash(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	hash := []byte("hello")
	err := repo.SetAccessTokenHash(userId, hash)
	assert.NoError(t, err)
	r, _ := repo.Find(userId)
	assert.True(t, bytes.Equal(r.HashPAT, hash))
}
