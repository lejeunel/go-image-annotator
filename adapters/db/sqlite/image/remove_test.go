package image

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnRemoveImageFromCollectionShouldFail(t *testing.T) {
	imRepo, clcRepo := MakeRepos()
	collectionId := clc.NewCollectionId()
	clcRepo.Create(clc.NewCollection(collectionId, "a-collection"))
	imageId := im.NewImageId()
	imRepo.AddImage(imageId, nil, im.ImageSpecs{})

	imRepo.AddToCollection(imageId, collectionId)
	imRepo.Db.Close()
	err := imRepo.RemoveImageFromCollection(imageId, collectionId)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRemoveImageFromCollection(t *testing.T) {
	imRepo, clcRepo := MakeRepos()
	collectionId := clc.NewCollectionId()
	clcRepo.Create(clc.NewCollection(collectionId, "a-collection"))
	imageId := im.NewImageId()
	imRepo.AddImage(imageId, nil, im.ImageSpecs{})

	imRepo.AddToCollection(imageId, collectionId)
	err := imRepo.RemoveImageFromCollection(imageId, collectionId)
	assert.NoError(t, err)
	exists, _ := imRepo.ImageExistsInCollection(imageId, collectionId)
	assert.False(t, exists)
}
