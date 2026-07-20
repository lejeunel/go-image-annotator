package collection

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnCollectionCountShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteCollectionRepo(db)
	db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCountCollections(t *testing.T) {
	repo := NewSQLiteCollectionRepo(s.NewInMemory())
	CreateCollection(repo, "a-collection")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnCollectionListShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteCollectionRepo(db)
	db.Close()
	_, err := repo.List(pa.PaginationParams{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestListCollections(t *testing.T) {
	repo := NewSQLiteCollectionRepo(s.NewInMemory())
	CreateCollection(repo, "a-collection")
	CreateCollection(repo, "another-collection")
	cs, err := repo.List(pa.PaginationParams{Page: 1, PageSize: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cs))
	assert.False(t, cs[0].Name == cs[1].Name)
	assert.False(t, cs[0].CreatedAt.Equal(cs[1].CreatedAt))
}
