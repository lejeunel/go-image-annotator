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
	imRepo.AddImage(imageId, nil, im.Specs{})

	imRepo.AddToCollection(imageId, collectionId)
	db.Close()
	err := imRepo.RemoveImageFromCollection(imageId, collectionId)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRemoveImageFromCollection(t *testing.T) {
	imRepo, clcRepo := MakeRepos(s.NewInMemory())
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	clcRepo.Create(collection)
	imageId := im.NewImageId()
	imRepo.AddImage(imageId, nil, im.Specs{})

	imRepo.AddToCollection(imageId, collection.Id)
	err := imRepo.RemoveImageFromCollection(imageId, collection.Id)
	assert.NoError(t, err)
	exists, _ := imRepo.ImageExistsInCollection(imageId, collection.Name)
	assert.False(t, exists)
}
