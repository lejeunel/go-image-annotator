package label

import (
	"errors"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestRetrieveNonExistingShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	CreateLabel(repo, "a-label")
	_, err := repo.FindLabel("non-existing-label")
	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestInternalErrOnFindShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	CreateLabel(repo, "a-label")
	repo.Db.Close()
	_, err := repo.FindLabel("a-label")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestRetrieve(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	label, _ := CreateLabel(repo, "a-label")
	r, err := repo.FindLabel("a-label")
	if err != nil {
		t.Fatalf("expected no error on find, got %v", err)
	}
	if (r.Name != label.Name) || (r.Description != label.Description) || (r.Id != label.Id) {
		t.Fatalf("expected to retrieve name %v, description %v, and id %v, got %v, %v, %v",
			label.Name, label.Description, label.Id, r.Name, r.Description, r.Id)

	}

}
