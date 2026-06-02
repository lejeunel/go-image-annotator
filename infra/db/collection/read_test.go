package collection

import (
	"errors"
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
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	CreateCollection(repo, "a-collection")
	repo.Db.Close()
	_, err := repo.FindCollectionByName("a-collection")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestRetrieve(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithDescription("a-description"), clc.WithCreatedAt(time.Now()), clc.WithGroup("a-group"))
	repo.Create(c)
	r, err := repo.FindCollectionByName("a-collection")
	assert.NoError(t, err, "expected no error on find")
	assert.Equal(t, c.Name, r.Name)
	assert.Equal(t, c.Description, r.Description)
	assert.Equal(t, c.Id, r.Id)
	assert.Equal(t, c.Group, r.Group)
	assert.Equal(t, r.CreatedAt.Equal(c.CreatedAt), true)

}

func TestRetrieveGroupOfCollection(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	c := clc.NewCollection(clc.NewCollectionId(), "a-collection",
		clc.WithDescription("a-description"), clc.WithCreatedAt(time.Now()), clc.WithGroup("a-group"))
	repo.Create(c)
	group, err := repo.GroupOfCollection("a-collection")
	assert.NoError(t, err, "expected no error on find")
	assert.Equal(t, c.Group, *group)

}
