package image

import (
	s "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	sc "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos/collection"
	ist "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ImageListingTestingRepos struct {
	Image      SQLiteImageRepo
	Collection sc.SQLiteCollectionRepo
}

func CreateSingleImageCollection(repos ImageListingTestingRepos, collectionName string) (im.Image, clc.Collection) {
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(collection)
	imageId := im.NewImageId()
	image := im.NewImage(imageId, collection)
	repos.Image.AddImage(image.Id, nil, im.ImageSpecs{})
	repos.Image.AddToCollection(image.Id, collection.Id)
	return image, collection
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
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestListOneImage(t *testing.T) {
	repos := NewImageListingTestRepos()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(collection)
	image := im.NewImage(im.NewImageId(), collection)
	repos.Image.AddImage(image.Id, nil, im.ImageSpecs{})
	repos.Image.AddToCollection(image.Id, collection.Id)

	r, _ := repos.Image.List(ist.FilteringParams{PageSize: 2, Page: 1})
	assert.Equal(t, 1, len(*r))
}

func TestListOneImageInGivenCollection(t *testing.T) {
	repos := NewImageListingTestRepos()

	firstImage, firstCollection := CreateSingleImageCollection(repos, "first-collection")
	CreateSingleImageCollection(repos, "second-collection")

	r, _ := repos.Image.List(ist.FilteringParams{Collection: &firstCollection.Name, PageSize: 2, Page: 1})
	assert.Equal(t, 1, len(*r))
	images := *r
	assert.True(t, images[0].ImageId == firstImage.Id.String())
	assert.True(t, images[0].Collection == firstCollection.Name)
}

func CreateImageInCollectionFromString(repo SQLiteImageRepo, collection clc.Collection, imageId string) im.Image {
	id, _ := im.NewImageIdFromString(imageId)
	image := im.NewImage(id, collection)
	repo.AddImage(image.Id, []byte(image.Id.String()), im.ImageSpecs{})
	repo.AddToCollection(image.Id, collection.Id)
	return image

}

func TestListImagesShouldBeOrderedById(t *testing.T) {
	repos := NewImageListingTestRepos()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(collection)
	CreateImageInCollectionFromString(repos.Image, collection, "11111111-1111-1111-1111-111111111111")
	image0 := CreateImageInCollectionFromString(repos.Image, collection, "00000000-0000-0000-0000-000000000000")

	r, _ := repos.Image.List(ist.FilteringParams{PageSize: 2, Page: 1})
	got := (*r)[0].ImageId
	assert.Equal(t, image0.Id.String(), got)
}
