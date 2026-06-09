package image

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalErrOnImageIsInCollectionShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.ImageExistsInCollection(im.NewImageId(), clc.NewCollectionId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddedImageToCollectionExists(t *testing.T) {
	repos := NewImageTestRepos()
	imageId, collectionId, _ := AddToCollection(repos, "a-collection", "the-hash")
	isAdded, err := repos.Image.ImageExistsInCollection(*imageId, *collectionId)
	assert.NoError(t, err)
	assert.True(t, isAdded)
}

func TestInternalErrOnImageExistsShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.ImageExists(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddedImageExists(t *testing.T) {
	repos := NewImageTestRepos()
	imageId, _, _ := AddToCollection(repos, "a-collection", "the-hash")
	exists, err := repos.Image.ImageExists(*imageId)
	assert.NoError(t, err)
	assert.True(t, exists)
}
