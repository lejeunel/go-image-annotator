package image

import (
	"errors"
	ist "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	s "github.com/lejeunel/go-image-annotator/infra/db"
	sc "github.com/lejeunel/go-image-annotator/infra/db/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

type ImageListingTestingRepos struct {
	Image      SQLiteImageRepo
	Collection sc.SQLiteCollectionRepo
}

func CreateSingleImageCollection(repos ImageListingTestingRepos, collectionName string) (im.Image, clc.Collection) {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(*collection)
	imageId := im.NewImageId()
	image := im.NewImage(imageId, *collection)
	repos.Image.AddImage(image.Id, nil, im.ImageSpecs{})
	repos.Image.AddToCollection(image.Id, collection.Id)
	return *image, *collection
}

func NewImageListingTestRepos() ImageListingTestingRepos {
	db := s.NewSQLiteDB(":memory:")
	return ImageListingTestingRepos{Image: *NewSQLiteImageRepo(db),
		Collection: *sc.NewSQLiteCollectionRepo(db),
	}
}

func TestInternalErrOnImageListShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.List(ist.FilteringParams{})
	if !errors.Is(err, e.ErrInternal) {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestListOneImage(t *testing.T) {
	repos := NewImageListingTestRepos()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(*collection)
	image := im.NewImage(im.NewImageId(), *collection)
	repos.Image.AddImage(image.Id, nil, im.ImageSpecs{})
	repos.Image.AddToCollection(image.Id, collection.Id)

	r, _ := repos.Image.List(ist.FilteringParams{PageSize: 2, Page: 1})
	if len(*r) != 1 {
		t.Fatalf("expected to retrieve one image, got %v", len(*r))
	}
}

func TestListOneImageInGivenCollection(t *testing.T) {
	repos := NewImageListingTestRepos()

	firstImage, firstCollection := CreateSingleImageCollection(repos, "first-collection")
	CreateSingleImageCollection(repos, "second-collection")

	r, _ := repos.Image.List(ist.FilteringParams{Collection: &firstCollection.Name, PageSize: 2, Page: 1})
	if len(*r) != 1 {
		t.Fatalf("expected to retrieve one image, got %v", len(*r))
	}
	images := *r
	if images[0].ImageId != firstImage.Id {
		t.Fatalf("expected to retrieve first image with id %v, got %v", firstImage.Id, images[0].ImageId)
	}
	if images[0].Collection != firstCollection.Name {
		t.Fatalf("expected that image belongs to collection named %v, got %v", firstCollection.Name, images[0].Collection)
	}
}

func CreateImageInCollectionFromString(repo SQLiteImageRepo, collection *clc.Collection, imageId string) *im.Image {
	id, _ := im.NewImageIdFromString(imageId)
	image := im.NewImage(id, *collection)
	repo.AddImage(image.Id, []byte(image.Id.String()), im.ImageSpecs{})
	repo.AddToCollection(image.Id, collection.Id)
	return image

}

func TestListImagesShouldBeOrderedById(t *testing.T) {
	repos := NewImageListingTestRepos()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(*collection)
	CreateImageInCollectionFromString(repos.Image, collection, "11111111-1111-1111-1111-111111111111")
	image0 := CreateImageInCollectionFromString(repos.Image, collection, "00000000-0000-0000-0000-000000000000")

	r, _ := repos.Image.List(ist.FilteringParams{PageSize: 2, Page: 1})
	got := (*r)[0].ImageId
	if got != image0.Id {
		t.Fatalf("expected to retrieve image with first id %v, got %v",
			image0.Id, got)
	}

}
