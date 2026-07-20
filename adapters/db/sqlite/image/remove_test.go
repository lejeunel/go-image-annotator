package image

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnRemoveImageFromCollectionShouldFail(t *testing.T) {
	db := s.NewInMemory()
	imRepo, clcRepo := MakeRepos(db)
	collectionId := clc.NewCollectionId()
	clcRepo.Create(clc.NewCollection(collectionId, "a-collection"))
	imageId := im.NewImageId()
	imRepo.AddImage(imageId, nil, im.ImageSpecs{})

	imRepo.AddToCollection(imageId, collectionId)
	db.Close()
	err := imRepo.RemoveImageFromCollection(imageId, collectionId)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRemoveImageFromCollection(t *testing.T) {
	imRepo, clcRepo := MakeRepos(s.NewInMemory())
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
