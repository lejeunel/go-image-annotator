package user

import (
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func CreateUser(repo SQLiteUserRepo, id string) (*u.User, error) {
	user := u.NewUser(id)
	if err := repo.Create(user); err != nil {
		return nil, err
	}
	return &user, nil

}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	repo.Db.Close()
	_, err := CreateUser(repo, userId)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreateAddsCount(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	_, err := CreateUser(repo, userId)
	assert.NoError(t, err)
	count, err := repo.Count()
	assert.Equal(t, int64(1), count)
	assert.NoError(t, err)
}

func TestNoCreatedUserDoNotExist(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	exists, err := repo.Exists(userId)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCreateAdmin(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	user := u.NewUser("admin", u.WithAdmin(true))
	err := repo.Create(user)
	assert.NoError(t, err)
	r, err := repo.Find("admin")
	assert.NoError(t, err)
	assert.Equal(t, true, r.IsAdmin)
}

func TestSetAdmin(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	user := u.NewUser("admin")
	repo.Create(user)
	err := repo.SetAdmin(user.Id, true)
	assert.NoError(t, err)
	r, _ := repo.Find(user.Id)
	assert.True(t, r.IsAdmin)
}

func TestPersonalAccessTokenHash(t *testing.T) {
	hash := []byte("pat-hash")
	repo := NewTestSQLiteUserRepo()
	user := u.NewUser("user@example.com",
		u.WithHashedPersonalAccessToken(hash))
	repo.Create(user)
	r, err := repo.Find("user@example.com")
	assert.NoError(t, err)
	assert.Equal(t, hash, r.HashPAT)
}

func TestPasswordHash(t *testing.T) {
	hash := []byte("password-hash")
	repo := NewTestSQLiteUserRepo()
	user := u.NewUser("user@example.com",
		u.WithHashedPassword(hash))
	repo.Create(user)
	r, err := repo.Find("user@example.com")
	assert.NoError(t, err)
	assert.Equal(t, hash, r.HashPassword)
}
