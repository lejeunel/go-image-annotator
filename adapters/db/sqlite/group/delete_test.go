package group

import (
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	clcr "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
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

func TestGroupPopulatedWithCollection(t *testing.T) {
	db := s.NewSQLiteDB(":memory:")
	clcRepo := clcr.NewSQLiteCollectionRepo(db)
	groupRepo := NewSQLiteGroupRepo(db)
	group, _ := CreateGroup(groupRepo, "a-group")
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection", clc.WithGroup(*group))
	clcRepo.Create(collection)
	isPopulated, err := groupRepo.IsPopulated(group.Name)
	assert.NoError(t, err)
	assert.True(t, *isPopulated)
}
