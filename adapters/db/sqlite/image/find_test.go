package image

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveImageIdByHash(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	imageId := im.NewImageId()
	hash := []byte("the-hash")
	err := repo.AddImage(imageId, hash, im.ImageSpecs{})
	assert.NoError(t, err)

	existingId, err := repo.FindImageIdByHash(hash)
	assert.NoError(t, err)
	assert.Equal(t, *existingId, imageId)
}

func TestRetrieveImageIdByNonExistingHashShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	imageId := im.NewImageId()
	repo.AddImage(imageId, nil, im.ImageSpecs{})
	_, err := repo.FindImageIdByHash([]byte("non-existing-hash"))
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestRetrieveImageIdByHashInternalErrShouldFail(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	_, err := repo.FindImageIdByHash(nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}
