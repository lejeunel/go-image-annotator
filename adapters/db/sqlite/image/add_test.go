package image

import (
	"testing"
	"time"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestAddSpecs(t *testing.T) {
	imRepo, _ := MakeRepos()
	id := im.NewImageId()

	specs := im.ImageSpecs{MIMEType: "the-mimetype", Width: 15, Height: 10}
	imRepo.AddImage(id, nil, specs)
	r, err := imRepo.GetSpecs(id)
	assert.NoError(t, err)
	assert.Equal(t, r.MIMEType, specs.MIMEType)
}

func TestCountAddedImageToCollection(t *testing.T) {
	imRepo, clcRepo := MakeRepos()
	collection := "a-collection"
	AddToCollection(imRepo, clcRepo, collection, "")
	count, err := imRepo.Count(im.CountingParams{Collection: &collection})
	assert.NoError(t, err)
	assert.Equal(t, 1, int(*count))
}

func TestCountAllImagesWhenAddingImageToCollection(t *testing.T) {
	imRepo, clcRepo := MakeRepos()
	AddToCollection(imRepo, clcRepo, "a-collection", "")
	count, err := imRepo.Count(im.CountingParams{})
	assert.NoError(t, err)
	assert.Equal(t, 1, int(*count))
}

func TestInternalErrOnCreateShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	err := repo.AddToCollection(im.NewImageId(), clc.NewCollectionId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestInternalErrOnIsCollectionPopulatedShouldFail(t *testing.T) {
	imRepo, clcRepo := MakeRepos()
	collectionName := "a-collection"
	AddToCollection(imRepo, clcRepo, collectionName, "the-hash")
	imRepo.Db.Close()
	_, err := clcRepo.IsPopulated(collectionName)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestIsCollectionPopulated(t *testing.T) {
	imRepo, clcRepo := MakeRepos()
	collectionName := "a-collection"
	AddToCollection(imRepo, clcRepo, collectionName, "the-hash")
	isPopulated, err := clcRepo.IsPopulated(collectionName)
	assert.NoError(t, err)
	assert.True(t, *isPopulated)
}

func TestCreatedAt(t *testing.T) {
	imRepo, clcRepo := MakeRepos()
	collectionName := "a-collection"
	now := time.Now()
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	clcRepo.Create(collection)
	imageId := im.NewImageId()
	err := imRepo.AddImage(imageId, nil, im.ImageSpecs{IngestedAt: now})
	assert.NoError(t, err)
	specs, err := imRepo.GetSpecs(imageId)
	assert.NoError(t, err)
	assert.Equal(t, now.Round(0), specs.IngestedAt.Round(0))
}
