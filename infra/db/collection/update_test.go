package collection

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	u "github.com/lejeunel/go-image-annotator/use-cases/collection/update"
	"testing"
)

func TestInternalErrOnCollectionUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	err := repo.Update(u.Model{})
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestUpdateCollection(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	collection, _ := CreateCollection(repo, "a-collection")
	newName := "new-collection-name"
	newDesc := "new-description"
	err := repo.Update(u.Model{Name: collection.Name, NewName: newName, NewDescription: newDesc})
	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	r, err := repo.FindCollectionByName(newName)
	if err != nil {
		t.Fatalf("expected to retrieve updated, got %v", err)
	}
	if (r.Name != newName) || (r.Description != newDesc) {
		t.Fatalf("expected to updated fields to name %v and description %v, got %v and %v",
			newName, newDesc, r.Name, r.Description)
	}
}
