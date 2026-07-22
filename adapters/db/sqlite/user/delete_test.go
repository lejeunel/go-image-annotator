package user

import (
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	repo.Db.Close()
	err := repo.Delete("user@example.com")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteLabel(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	user, _ := CreateUser(repo, "user@example.com")
	err := repo.Delete(user.Id)
	assert.NoError(t, err)
}
