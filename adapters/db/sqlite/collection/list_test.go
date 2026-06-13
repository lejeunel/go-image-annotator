package collection

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	l "github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnCollectionCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	_, err := repo.Count()
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCountCollections(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	CreateCollection(repo, "a-collection")
	count, _ := repo.Count()
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnCollectionListShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	_, err := repo.List(l.Request{})
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestListCollections(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	CreateCollection(repo, "a-collection")
	CreateCollection(repo, "another-collection")
	cs, err := repo.List(l.Request{Page: 1, PageSize: 2})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cs))
	assert.False(t, cs[0].Name == cs[1].Name)
	assert.False(t, cs[0].CreatedAt.Equal(cs[1].CreatedAt))
}
