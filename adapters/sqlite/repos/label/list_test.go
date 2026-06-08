package label

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	l "github.com/lejeunel/go-image-annotator/use-cases/label/list"
	"testing"
)

func TestInternalErrOnLabelCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.Count()
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestCountLabels(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	CreateLabel(repo, "a-label")
	count, _ := repo.Count()
	if count != 1 {
		t.Fatalf("expected label count %v, got %v", 1, count)
	}
}

func TestInternalErrOnLabelListShouldFail(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	repo.Db.Close()
	_, err := repo.List(l.Request{})
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestListLabels(t *testing.T) {
	repo := NewTestSQLiteLabelRepo()
	firstLabel, _ := CreateLabel(repo, "a-label")
	secondLabel, _ := CreateLabel(repo, "another-label")
	labels, err := repo.List(l.Request{Page: 1, PageSize: 2})
	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected two labels, got %v", len(labels))
	}
	if labels[0].Name == labels[1].Name {
		t.Fatalf("expected to retrieve two distinct labels with name %v and %v, got %v and %v",
			firstLabel.Name, secondLabel.Name, labels[0].Name, labels[1].Name)

	}
	if labels[0].Description != firstLabel.Description {
		t.Fatalf("expected to retrieve label with description %v , got %v",
			firstLabel.Description, labels[0].Description)
	}
}
