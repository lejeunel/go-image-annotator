package group

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	l "github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"testing"
)

func TestInternalErrOnCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := repo.Count()
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestCount(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	CreateGroup(repo, "a-group")
	count, _ := repo.Count()
	if *count != 1 {
		t.Fatalf("expected label count %v, got %v", 1, count)
	}
}

func TestInternalErrOnListShouldFail(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	repo.Db.Close()
	_, err := repo.List(l.Request{})
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestList(t *testing.T) {
	repo := NewTestSQLiteGroupRepo()
	first, _ := CreateGroup(repo, "a-group")
	second, _ := CreateGroup(repo, "another-group")
	cs, err := repo.List(l.Request{Page: 1, PageSize: 2})
	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	if len(cs) != 2 {
		t.Fatalf("expected two groups, got %v", len(cs))
	}
	if cs[0].Name == cs[1].Name {
		t.Fatalf("expected to retrieve two distinct groups with name %v and %v, got %v and %v",
			first.Name, second.Name, cs[0].Name, cs[1].Name)

	}
	if cs[0].Description != first.Description {
		t.Fatalf("expected to retrieve group with description %v , got %v",
			first.Description, cs[0].Description)
	}
}
