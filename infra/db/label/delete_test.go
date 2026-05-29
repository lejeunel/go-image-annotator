package label

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestCreatedLabelExists(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	label, _ := CreateLabel(repo, "a-label")
	exists, _ := repo.Exists(label.Name)
	if !exists {
		t.Fatal("expected that created label exists")
	}
}

func TestNonExistingLabelDoesNotExists(t *testing.T) {
	exists, _ := NewTestSQLiteLabelRepo().Exists("non-existing-label")
	if exists {
		t.Fatal("expected that non-existing label does not exist")
	}
}

func TestInternalErrOnLabelExistsShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.Exists("")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestInternalErrOnDeleteShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	err := repo.Delete("a-label")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestDeleteLabel(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	label, _ := CreateLabel(repo, "a-label")
	err := repo.Delete(label.Name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
