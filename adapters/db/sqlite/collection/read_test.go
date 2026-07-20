package collection

import (
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewSQLiteCollectionRepo(s.NewInMemory())
	CreateCollection(repo, "a-collection")
	_, err := repo.FindCollectionByName("non-existing-collection")
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteCollectionRepo(db)
	CreateCollection(repo, "a-collection")
	db.Close()
	_, err := repo.FindCollectionByName("a-collection")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRetrieve(t *testing.T) {
	repo := NewSQLiteCollectionRepo(s.NewInMemory())
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
