package collection

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreatedCollectionExists(t *testing.T) {
	repo := NewSQLiteCollectionRepo(s.NewInMemory())
	collection, _ := CreateCollection(repo, "a-collection")
	exists, _ := repo.Exists(collection.Name)
	assert.True(t, exists)
}

func TestNonExistingCollectionDoesNotExists(t *testing.T) {
	exists, _ := NewSQLiteCollectionRepo(s.NewInMemory()).Exists("non-existing-collection")
	assert.False(t, exists)
}

func TestInternalErrOnCollectionExistsShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteCollectionRepo(db)
	db.Close()
	_, err := repo.Exists("")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteCollectionRepo(db)
	db.Close()
	err := repo.Delete("a-collection")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteCollection(t *testing.T) {
	repo := NewSQLiteCollectionRepo(s.NewInMemory())
	collection, _ := CreateCollection(repo, "a-collection")
	err := repo.Delete(collection.Name)
	assert.NoError(t, err)
}
