package user

import (
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func CreateUser(repo *SQLiteUserRepo, id string) (*u.User, error) {
	user := u.NewUser(id)
	if err := repo.Create(user); err != nil {
		return nil, err
	}
	return &user, nil

}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	repo.Db.Close()
	_, err := CreateUser(repo, "user@example.com")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreateAddsCount(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	_, err := CreateUser(repo, "user@example.com")
	assert.NoError(t, err)
	count, err := repo.Count()
	assert.Equal(t, int64(1), count)
	assert.NoError(t, err)
}
