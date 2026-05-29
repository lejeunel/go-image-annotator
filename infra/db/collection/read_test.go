package collection

import (
	"errors"
	"testing"
	"time"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
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
		clc.WithDescription("a-description"), clc.WithCreatedAt(time.Now()))
	repo.Create(*c)
	r, err := repo.FindCollectionByName("a-collection")
	if err != nil {
		t.Fatalf("expected no error on find, got %v", err)
	}
	if (r.Name != c.Name) || (r.Description != c.Description) || (r.Id != c.Id) || !r.CreatedAt.Equal(c.CreatedAt) {
		t.Fatalf("expected to retrieve name %v, description %v, id %v, created at %v, got %v, %v, %v, %v",
			c.Name, c.Description, c.Id, c.CreatedAt,
			r.Name, r.Description, r.Id, r.CreatedAt)

	}

}
