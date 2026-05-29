package label

import (
	"errors"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestInternalErrOnFetchAll(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.FetchAll()
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestFetchAll(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Create(*l.NewLabel(l.NewLabelId(), "first-label"))
	repo.Create(*l.NewLabel(l.NewLabelId(), "second-label"))
	labels, _ := repo.FetchAll()
	if len(labels) != 2 {
		t.Fatalf("expected to retrieve 2 labels, got %v", len(labels))
	}
}
