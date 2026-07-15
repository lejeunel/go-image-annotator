package image

import (
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	sc "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/collection"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
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
	return ImageListingTestingRepos{
		Image:      NewSQLiteImageRepo(db),
		Collection: sc.NewSQLiteCollectionRepo(db),
	}
}

func TestInternalErrOnImageListShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.Slice(im.FilteringParams{}, im.PaginationParams{}, im.OrderingParams{})
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

	r, _ := repos.Image.Slice(im.FilteringParams{}, im.PaginationParams{PageSize: 2, Page: 1}, im.OrderingParams{})
	assert.Equal(t, 1, len(r))
}

func TestListOneImageInGivenCollection(t *testing.T) {
	repos := NewImageListingTestRepos()

	firstImage, firstCollection := CreateSingleImageCollection(repos, "first-collection")
	CreateSingleImageCollection(repos, "second-collection")

	r, _ := repos.Image.Slice(im.FilteringParams{Collection: &firstCollection.Name}, im.PaginationParams{PageSize: 2, Page: 1}, im.OrderingParams{})
	assert.Equal(t, 1, len(r))
	images := r
	assert.True(t, images[0].ImageId == firstImage.Id)
	assert.True(t, images[0].Collection == firstCollection.Name)
}

func CreateImageInCollectionFromString(repo SQLiteImageRepo, collection clc.Collection, imageId string) im.Image {
	id, _ := im.NewImageIdFromString(imageId)
	image := im.NewImage(id, collection)
	repo.AddImage(image.Id, []byte(image.Id.String()), im.ImageSpecs{})
	repo.AddToCollection(image.Id, collection.Id)
	return image

}

func TestListImagesOrderedById(t *testing.T) {
	repos := NewImageListingTestRepos()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(collection)
	CreateImageInCollectionFromString(repos.Image, collection, "11111111-1111-1111-1111-111111111111")
	image0 := CreateImageInCollectionFromString(repos.Image, collection, "00000000-0000-0000-0000-000000000000")

	r, _ := repos.Image.Slice(im.FilteringParams{}, im.PaginationParams{PageSize: 2, Page: 1}, im.OrderingParams{})
	got := r[0].ImageId
	assert.Equal(t, image0.Id, got)
}

func TestListImagesOrderedByIngestTime(t *testing.T) {
	repos := NewImageListingTestRepos()
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	repos.Collection.Create(collection)
	firstId, _ := im.NewImageIdFromString("11111111-1111-1111-1111-111111111111")
	secondId, _ := im.NewImageIdFromString("00000000-0000-0000-0000-000000000000")
	firstImage := im.NewImage(firstId, collection)
	secondImage := im.NewImage(secondId, collection)
	repos.Image.AddImage(firstImage.Id, []byte("first-hash"), im.ImageSpecs{IngestedAt: time.Now()})
	repos.Image.AddImage(secondImage.Id, []byte("second-hash"), im.ImageSpecs{IngestedAt: time.Now()})
	repos.Image.AddToCollection(firstImage.Id, collection.Id)
	repos.Image.AddToCollection(secondImage.Id, collection.Id)

	r, err := repos.Image.Slice(im.FilteringParams{}, im.PaginationParams{PageSize: 2, Page: 1}, im.OrderingParams{IngestTime: true})
	assert.NoError(t, err)
	assert.Equal(t, r[0].ImageId, firstImage.Id)
	assert.Equal(t, r[1].ImageId, secondImage.Id)
}

func TestIterateImages(t *testing.T) {
	repos := NewImageListingTestRepos()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	repos.Collection.Create(collection)
	im0 := CreateImageInCollectionFromString(repos.Image, collection, "00000000-0000-0000-0000-000000000000")
	im1 := CreateImageInCollectionFromString(repos.Image, collection, "11111111-1111-1111-1111-111111111111")

	res := []im.BaseImage{}
	for batch, err := range repos.Image.Iterate(im.FilteringParams{}, 1) {
		assert.NoError(t, err)
		res = append(res, batch)
	}
	assert.Equal(t, 2, len(res))
	assert.Equal(t, im0.Id.String(), res[0].ImageId.String())
	assert.Equal(t, im1.Id.String(), res[1].ImageId.String())
}
