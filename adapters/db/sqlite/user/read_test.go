package user

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	_, err := repo.Find("user@example.com")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	CreateUser(repo, "user@example.com")
	repo.Db.Close()
	_, err := repo.Find("user@example.com")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRetrieve(t *testing.T) {
	repo := NewSQLiteUserRepo(s.NewInMemory())
	user, _ := CreateUser(repo, "user@example.com")
	r, err := repo.Find("user@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Id, r.Id)
}
