package image

import (
	"errors"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestInternalErrOnImageIsInCollectionShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.ImageExistsInCollection(im.NewImageId(), clc.NewCollectionId())
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestAddedImageToCollectionExists(t *testing.T) {
	repos := NewImageTestRepos()
	imageId, collectionId, _ := AddToCollection(repos, "a-collection", "the-hash")
	isAdded, err := repos.Image.ImageExistsInCollection(*imageId, *collectionId)
	if err != nil {
		t.Fatalf("expected no error when checking existence of image in collection, got %v", err)
	}
	if !isAdded {
		t.Fatal("expected that image is added to collection")
	}
}

func TestInternalErrOnImageExistsShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.ImageExists(im.NewImageId())
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestAddedImageExists(t *testing.T) {
	repos := NewImageTestRepos()
	imageId, _, _ := AddToCollection(repos, "a-collection", "the-hash")
	exists, err := repos.Image.ImageExists(*imageId)
	if err != nil {
		t.Fatalf("expected no error when checking existence of image, got %v", err)
	}
	if !exists {
		t.Fatal("expected that image exists")
	}
}
