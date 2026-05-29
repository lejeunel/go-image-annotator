package label

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	u "github.com/lejeunel/go-image-annotator/use-cases/label/update"
	"testing"
)

func TestInternalErrOnLabelUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	err := repo.Update(u.Model{})
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestUpdateLabel(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	label, _ := CreateLabel(repo, "a-label")
	newName := "new-label-name"
	newDesc := "new-description"
	err := repo.Update(u.Model{Name: label.Name, NewName: newName, NewDescription: newDesc})
	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	r, err := repo.FindLabelByName(newName)
	if err != nil {
		t.Fatalf("expected to retrieve updated, got %v", err)
	}
	if (r.Name != newName) || (r.Description != newDesc) {
		t.Fatalf("expected to updated fields to name %v and description %v, got %v and %v",
			newName, newDesc, r.Name, r.Description)
	}
}
