package user

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	repo.Db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCount(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, "user@example.com")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(count))
}

func TestInternalErrOnListShouldFail(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	repo.Db.Close()
	_, err := repo.List(list.Request{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestList(t *testing.T) {
	repo := NewTestSQLiteUserRepo()
	CreateUser(repo, "user@example.com")
	CreateUser(repo, "another-user@example.com")
	users, err := repo.List(list.Request{Page: 1, PageSize: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
	assert.False(t, users[0].Id == users[1].Id)
}
