package collection

import (
	"testing"
	"time"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	CreateCollection(repo, "a-collection")
	_, err := repo.FindCollectionByName("non-existing-collection")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	CreateCollection(repo, "a-collection")
	repo.Db.Close()
	_, err := repo.FindCollectionByName("a-collection")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRetrieve(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithDescription("a-description"),
		clc.WithCreatedAt(time.Now()))
	repo.Create(c)
	r, err := repo.FindCollectionByName("a-collection")
	assert.NoError(t, err, "expected no error on find")
	assert.Equal(t, c.Name, r.Name)
	assert.Equal(t, c.Description, r.Description)
	assert.Equal(t, c.Id, r.Id)
	assert.Equal(t, r.CreatedAt.Equal(c.CreatedAt), true)

}
