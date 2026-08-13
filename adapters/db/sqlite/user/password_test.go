package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePassword(t *testing.T) {
	repo := NewUserRepo(s.NewInMemory())
	CreateUser(repo, userId)
	pwHash := []byte("hello")
	err := repo.UpdatePassword(userId, pwHash)
	assert.NoError(t, err)
	r, _ := repo.Find(userId)
	assert.Equal(t, pwHash, r.HashPassword)
}
