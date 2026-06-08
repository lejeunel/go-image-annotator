package collection

import (
	"errors"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	l "github.com/lejeunel/go-image-annotator/use-cases/collection/list"
	"testing"
)

func TestInternalErrOnCollectionCountShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	_, err := repo.Count()
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestCountCollections(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	CreateCollection(repo, "a-collection")
	count, _ := repo.Count()
	if *count != 1 {
		t.Fatalf("expected label count %v, got %v", 1, count)
	}
}

func TestInternalErrOnCollectionListShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	_, err := repo.List(l.Request{})
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestListCollections(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	first, _ := CreateCollection(repo, "a-collection")
	second, _ := CreateCollection(repo, "another-collection")
	cs, err := repo.List(l.Request{Page: 1, PageSize: 2})
	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	if len(cs) != 2 {
		t.Fatalf("expected two collections, got %v", len(cs))
	}
	if cs[0].Name == cs[1].Name {
		t.Fatalf("expected to retrieve two distinct collections with name %v and %v, got %v and %v",
			first.Name, second.Name, cs[0].Name, cs[1].Name)

	}
	if cs[0].Description != first.Description {
		t.Fatalf("expected to retrieve collection with description %v , got %v",
			first.Description, cs[0].Description)
	}
	if !cs[0].CreatedAt.Equal(first.CreatedAt) {
		t.Fatalf("expected to retrieve created at %v , got %v",
			first.CreatedAt, cs[0].CreatedAt)
	}
}
