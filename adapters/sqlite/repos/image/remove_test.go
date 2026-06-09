package image

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnRemoveImageFromCollectionShouldFail(t *testing.T) {
	repos := NewImageTestRepos()
	collectionId := clc.NewCollectionId()
	repos.Collection.Create(clc.NewCollection(collectionId, "a-collection"))
	imageId := im.NewImageId()
	repos.Image.AddImage(imageId, nil, im.ImageSpecs{})

	repos.Image.AddToCollection(imageId, collectionId)
	repos.Image.Db.Close()
	err := repos.Image.RemoveImageFromCollection(imageId, collectionId)
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestRemoveImageFromCollection(t *testing.T) {
	repos := NewImageTestRepos()
	collectionId := clc.NewCollectionId()
	repos.Collection.Create(clc.NewCollection(collectionId, "a-collection"))
	imageId := im.NewImageId()
	repos.Image.AddImage(imageId, nil, im.ImageSpecs{})

	repos.Image.AddToCollection(imageId, collectionId)
	err := repos.Image.RemoveImageFromCollection(imageId, collectionId)
	assert.NoError(t, err)
	exists, _ := repos.Image.ImageExistsInCollection(imageId, collectionId)
	assert.False(t, exists)
}
