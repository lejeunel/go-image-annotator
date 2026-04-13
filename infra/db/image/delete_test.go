package image

import (
	"errors"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	"testing"
)

func TestHandleInternalErrOnDeleteImage(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	err := repo.Delete(im.NewImageId())
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestDeleteImage(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	id := im.NewImageId()
	repo.AddImage(id, nil, "")
	err := repo.Delete(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
