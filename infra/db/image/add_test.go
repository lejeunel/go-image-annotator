package image

import (
	"errors"
	ist "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	"testing"
)

func TestAddSpecs(t *testing.T) {
	repos := NewImageTestRepos()
	id := im.NewImageId()

	specs := im.ImageSpecs{MIMEType: "the-mimetype", Width: 15, Height: 10}
	repos.Image.AddImage(id, nil, specs)
	r, err := repos.Image.GetSpecs(id)
	if err != nil {
		t.Fatalf("expected no error when retrieving specs, got %v", err)
	}
	if r.MIMEType != specs.MIMEType {
		t.Fatalf("expected to retrieve mimetype %v, got %v", specs.MIMEType, r.MIMEType)
	}
}

func TestCountAddedImageToCollection(t *testing.T) {
	repos := NewImageTestRepos()
	collection := "a-collection"
	AddToCollection(repos, collection, "")
	count, err := repos.Image.Count(ist.CountingParams{Collection: &collection})
	if err != nil {
		t.Fatalf("expected no error when counting images in collection, got %v", err)
	}
	if *count != 1 {
		t.Fatalf("expected that one image is added to collection, got %v", *count)
	}
}

func TestCountAllImagesWhenAddingImageToCollection(t *testing.T) {
	repos := NewImageTestRepos()
	AddToCollection(repos, "a-collection", "")
	count, err := repos.Image.Count(ist.CountingParams{})
	if err != nil {
		t.Fatalf("expected no error when counting images in collection, got %v", err)
	}
	if *count != 1 {
		t.Fatalf("expected that one image is added to collection, got %v", *count)
	}
}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	err := repo.AddToCollection(im.NewImageId(), clc.NewCollectionId())
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestInternalErrOnIsCollectionPopulatedShouldFail(t *testing.T) {
	repos := NewImageTestRepos()
	collectionName := "a-collection"
	AddToCollection(repos, collectionName, "the-hash")
	repos.Image.Db.Close()
	_, err := repos.Collection.IsPopulated(collectionName)
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestIsCollectionPopulated(t *testing.T) {
	repos := NewImageTestRepos()
	collectionName := "a-collection"
	AddToCollection(repos, collectionName, "the-hash")
	isPopulated, err := repos.Collection.IsPopulated(collectionName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !(*isPopulated) {
		t.Fatal("expected populated collection, got")
	}
}
