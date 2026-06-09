package collection

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatedCollectionExists(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	collection, _ := CreateCollection(repo, "a-collection")
	exists, _ := repo.Exists(collection.Name)
	assert.True(t, exists)
}

func TestNonExistingCollectionDoesNotExists(t *testing.T) {
	exists, _ := NewTestSQLiteCollectionRepo().Exists("non-existing-collection")
	assert.False(t, exists)
}

func TestInternalErrOnCollectionExistsShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	_, err := repo.Exists("")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	err := repo.Delete("a-collection")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestDeleteCollection(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	collection, _ := CreateCollection(repo, "a-collection")
	err := repo.Delete(collection.Name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
