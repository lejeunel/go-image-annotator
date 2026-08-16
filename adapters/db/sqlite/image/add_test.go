package image

import (
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestAddSpecs(t *testing.T) {
	imRepo, _, _ := SetupAdd(s.NewInMemory())
	id := im.NewImageId()

	specs := im.Specs{MIMEType: "the-mimetype", Width: 15, Height: 10}
	imRepo.AddImage(id, nil, specs)
	r, err := imRepo.GetSpecs(id)
	assert.NoError(t, err)
	assert.Equal(t, r.MIMEType, specs.MIMEType)
}

func TestCountAddedImageToCollection(t *testing.T) {
	imRepo, clcRepo, _ := SetupAdd(s.NewInMemory())
	collection := "a-collection"
	imageId, _, _ := AddToCollection(imRepo, clcRepo, collection, "")
	count, err := imRepo.Count("collection=\"a-collection\"")
	assert.NoError(t, err)
	assert.Equal(t, 1, int(*count))
	isUsed, err := imRepo.IsUsed(*imageId)
	assert.NoError(t, err)
	assert.True(t, *isUsed)
}

func TestCountAllImagesWhenAddingImageToCollection(t *testing.T) {
	imRepo, clcRepo, _ := SetupAdd(s.NewInMemory())
	AddToCollection(imRepo, clcRepo, "a-collection", "")
	count, err := imRepo.Count("collection=\"a-collection\"")
	assert.NoError(t, err)
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	imr, _, db := SetupAdd(s.NewInMemory())
	db.Close()
	err := imr.AddToCollection(im.NewImageId(), "a-collection")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnIsCollectionPopulatedShouldFail(t *testing.T) {
	db := s.NewInMemory()
	imRepo, clcRepo, _ := SetupAdd(db)
	collectionName := "a-collection"
	AddToCollection(imRepo, clcRepo, collectionName, "the-hash")
	db.Close()
	_, err := clcRepo.IsPopulated(collectionName)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestIsCollectionPopulated(t *testing.T) {
	db := s.NewInMemory()
	imRepo, clcRepo, _ := SetupAdd(db)
	collectionName := "a-collection"
	AddToCollection(imRepo, clcRepo, collectionName, "the-hash")
	isPopulated, err := clcRepo.IsPopulated(collectionName)
	assert.NoError(t, err)
	assert.True(t, *isPopulated)
}

func TestCreatedAt(t *testing.T) {
	db := s.NewInMemory()
	imRepo, clcRepo, _ := SetupAdd(db)
	collectionName := "a-collection"
	now := time.Now()
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	clcRepo.Create(collection)
	imageId := im.NewImageId()
	err := imRepo.AddImage(imageId, nil, im.Specs{IngestedAt: now})
	assert.NoError(t, err)
	specs, err := imRepo.GetSpecs(imageId)
	assert.NoError(t, err)
	assert.True(t, now.Equal(specs.IngestedAt))
}
