package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	repo := NewUserRepo(s.NewInMemory())
	repo.Db.Close()
	err := repo.Delete("user@example.com")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteUser(t *testing.T) {
	repo := NewUserRepo(s.NewInMemory())
	user, _ := CreateUser(repo, "user@example.com")
	err := repo.Delete(user.Id)
	assert.NoError(t, err)
}
