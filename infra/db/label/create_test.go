package label

import (
	"errors"
	"testing"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func CreateLabel(repo *SQLiteLabelRepo, name string) (*lbl.Label, error) {
	label := lbl.NewLabel(lbl.NewLabelId(), name, lbl.WithDescription("a-description"))
	if err := repo.Create(*label); err != nil {
		return nil, err
	}
	return label, nil

}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := CreateLabel(repo, "a-label")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestCreateAddsCount(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	_, err := CreateLabel(repo, "a-label")
	if err != nil {
		t.Fatalf("expected no error on create but got %v", err)
	}
	count, err := repo.Count()
	if err != nil {
		t.Fatalf("expected no error on count, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count of 1, got %v", count)
	}

}
