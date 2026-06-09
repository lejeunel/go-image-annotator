package group

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	u "github.com/lejeunel/go-image-annotator/use-cases/group/update"
	"testing"
)

func TestInternalErrOnUpdateShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	err := repo.Update(u.Model{})
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	group, _ := CreateGroup(repo, "a-group")
	newName := "new-group-name"
	newDesc := "new-description"
	err := repo.Update(u.Model{Name: group.Name, NewName: newName, NewDescription: newDesc})
	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	r, err := repo.Find(newName)
	if err != nil {
		t.Fatalf("expected to retrieve updated, got %v", err)
	}
	if (r.Name != newName) || (r.Description != newDesc) {
		t.Fatalf("expected to updated fields to name %v and description %v, got %v and %v",
			newName, newDesc, r.Name, r.Description)
	}
}
