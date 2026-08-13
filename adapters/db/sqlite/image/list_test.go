package image

import (
	"testing"
	"time"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnImageListShouldFail(t *testing.T) {
	imr, _, db := SetupList()
	db.Close()
	_, err := imr.Slice("collection:\"my-collection\"", pa.PaginationParams{}, "")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestListOneImage(t *testing.T) {
	imr, cr, _ := SetupList()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	cr.Create(collection)
	image := im.NewImage(im.NewImageId(), collection)
	imr.AddImage(image.Id, nil, im.Specs{})
	imr.AddToCollection(image.Id, collection.Name)

	r, err := imr.Slice("", pa.PaginationParams{PageSize: 2, Page: 1}, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r))
}

func TestListOneImageInGivenCollection(t *testing.T) {
	imr, cr, _ := SetupList()
	firstImage, firstCollection := CreateSingleImageCollection(imr, cr, "first-collection")
	CreateSingleImageCollection(imr, cr, "second-collection")

	r, _ := imr.Slice("collection=first-collection", pa.PaginationParams{PageSize: 2, Page: 1}, "")
	assert.Equal(t, 1, len(r))
	images := r
	assert.True(t, images[0].ImageId == firstImage.Id)
	assert.True(t, images[0].Collection == firstCollection.Name)
}

func CreateImageInCollectionFromString(
	repo ImageRepo,
	collection clc.Collection,
	imageId string,
) im.Image {
	id, _ := im.NewImageIdFromString(imageId)
	image := im.NewImage(id, collection)
	repo.AddImage(image.Id, []byte(image.Id.String()), im.Specs{})
	repo.AddToCollection(image.Id, collection.Name)
	return image
}

func TestListImagesOrderedById(t *testing.T) {
	imr, cr, _ := SetupList()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	cr.Create(collection)
	CreateImageInCollectionFromString(
		imr,
		collection,
		st.FakeUUIDFromInt(1),
	)
	image0 := CreateImageInCollectionFromString(
		imr,
		collection,
		st.FakeUUIDFromInt(0),
	)

	r, _ := imr.Slice(
		"", pa.PaginationParams{PageSize: 2, Page: 1}, "")
	got := r[0].ImageId
	assert.Equal(t, image0.Id, got)
}

func TestListImagesOrderedByIngestTime(t *testing.T) {
	imr, cr, _ := SetupList()
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	cr.Create(collection)
	firstId, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(0))
	secondId, _ := im.NewImageIdFromString(st.FakeUUIDFromInt(1))
	firstImage := im.NewImage(firstId, collection)
	secondImage := im.NewImage(secondId, collection)
	imr.AddImage(firstImage.Id, []byte("first-hash"), im.Specs{IngestedAt: time.Now()})
	imr.AddImage(secondImage.Id, []byte("second-hash"), im.Specs{IngestedAt: time.Now()})
	imr.AddToCollection(firstImage.Id, collection.Name)
	imr.AddToCollection(secondImage.Id, collection.Name)

	r, err := imr.Slice(
		"",
		pa.PaginationParams{PageSize: 2, Page: 1},
		"ingested_at:asc")
	assert.NoError(t, err)
	assert.Equal(t, r[0].ImageId, firstImage.Id)
	assert.Equal(t, r[1].ImageId, secondImage.Id)
}

func TestIterateImages(t *testing.T) {
	imr, cr, _ := SetupList()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	cr.Create(collection)
	im0 := CreateImageInCollectionFromString(
		imr,
		collection,
		st.FakeUUIDFromInt(0),
	)
	im1 := CreateImageInCollectionFromString(
		imr,
		collection,
		st.FakeUUIDFromInt(1),
	)

	res := []im.BaseImage{}
	for batch, err := range imr.Iterate("", 1) {
		assert.NoError(t, err)
		res = append(res, batch)
	}
	assert.Equal(t, 2, len(res))
	assert.Equal(t, im0.Id.String(), res[0].ImageId.String())
	assert.Equal(t, im1.Id.String(), res[1].ImageId.String())
}
