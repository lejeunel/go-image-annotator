package group

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatedExists(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	group, _ := CreateGroup(repo, "a-group")
	exists, _ := repo.Exists(group.Name)
	assert.True(t, *exists)
}

func TestNonExistingDoesNotExists(t *testing.T) {
	exists, _ := NewTestSQLiteGroupRepo().Exists("non-existing-group")
	assert.False(t, *exists)
}

func TestInternalErrOnExistsShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := repo.Exists("")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	err := repo.Delete("a-group")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDelete(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	grp, _ := CreateGroup(repo, "a-group")
	err := repo.Delete(grp.Name)
	assert.Nil(t, err)
}
