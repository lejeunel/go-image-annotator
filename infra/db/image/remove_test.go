package image

import (
	"errors"
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestInternalErrOnRemoveImageFromCollectionShouldFail(t *testing.T) {
	repos := NewImageTestRepos()
	collectionId := clc.NewCollectionId()
	repos.Collection.Create(*clc.NewCollection(collectionId, "a-collection"))
	imageId := im.NewImageId()
	repos.Image.AddImage(imageId, nil, im.ImageSpecs{})

	repos.Image.AddToCollection(imageId, collectionId)
	repos.Image.Db.Close()
	err := repos.Image.RemoveImageFromCollection(imageId, collectionId)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestRemoveImageFromCollection(t *testing.T) {
	repos := NewImageTestRepos()
	collectionId := clc.NewCollectionId()
	repos.Collection.Create(*clc.NewCollection(collectionId, "a-collection"))
	imageId := im.NewImageId()
	repos.Image.AddImage(imageId, nil, im.ImageSpecs{})

	repos.Image.AddToCollection(imageId, collectionId)
	err := repos.Image.RemoveImageFromCollection(imageId, collectionId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	exists, _ := repos.Image.ImageExistsInCollection(imageId, collectionId)
	if exists {
		t.Fatal("expected that removed image does not exist")
	}
}
