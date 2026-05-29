package collection

import (
	"errors"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteCollectionRepo()
	repo.Db.Close()
	_, err := CreateCollection(repo, "a-collection")
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestCreate(t *testing.T) {
	_, err := CreateCollection(NewTestSQLiteCollectionRepo(), "a-collection")
	if err != nil {
		t.Fatalf("expected no error on create but got %v", err)
	}

}
