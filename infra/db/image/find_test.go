package image

import (
	"errors"
	"testing"

	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

func TestRetrieveImageIdByHash(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	imageId := im.NewImageId()
	hash := []byte("the-hash")
	err := repo.AddImage(imageId, hash, "")
	if err != nil {
		t.Fatalf("expected no error on adding image, got %v", err)
	}

	existingId, err := repo.FindImageIdByHash(hash)
	if err != nil {
		t.Fatalf("expected no error finding image by hash, got %v", err)
	}
	if *existingId != imageId {
		t.Fatalf("expected to retrieve image with identical hash and id %v, got %v", imageId, existingId)
	}
}

func TestRetrieveImageIdByNonExistingHashShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	imageId := im.NewImageId()
	repo.AddImage(imageId, nil, "")
	_, err := repo.FindImageIdByHash([]byte("non-existing-hash"))

	if !errors.Is(err, e.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestRetrieveImageIdByHashInternalErrShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.FindImageIdByHash(nil)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}
