package collection

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCollectionUpdateShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteCollectionRepo(db)
	db.Close()
	err := repo.Update(clc.UpdateModel{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestUpdateCollection(t *testing.T) {
	repo := NewSQLiteCollectionRepo(s.NewInMemory())
	collection, _ := CreateCollection(repo, "a-collection")
	newName := "new-collection-name"
	newDesc := "new-description"
	err := repo.Update(clc.UpdateModel{Name: collection.Name, NewName: newName, NewDescription: newDesc})
	assert.NoError(t, err)
	r, err := repo.FindCollectionByName(newName)
	assert.NoError(t, err)
	assert.Equal(t, newName, r.Name)
	assert.Equal(t, newDesc, r.Description)
}
