package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdatePassword(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, userId)
	pwHash := []byte("hello")
	err := repo.UpdatePassword(userId, pwHash)
	assert.NoError(t, err)
	r, _ := repo.Find(userId)
	assert.Equal(t, pwHash, r.HashPassword)
}
