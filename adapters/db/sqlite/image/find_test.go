package image

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestFindImageIdByHash(t *testing.T) {
	repo, _ := BaseSetup()
	imageId := im.NewImageId()
	hash := []byte("the-hash")
	err := repo.AddImage(imageId, hash, im.Specs{})
	assert.NoError(t, err)

	existingId, err := repo.FindImageIdByHash(hash)
	assert.NoError(t, err)
	assert.Equal(t, *existingId, imageId)
}

func TestFindImageIdByNonExistingHashShouldFail(t *testing.T) {
	repo, _ := BaseSetup()
	imageId := im.NewImageId()
	repo.AddImage(imageId, nil, im.Specs{})
	_, err := repo.FindImageIdByHash([]byte("non-existing-hash"))
	assert.ErrorIs(t, err, e.ErrNotFound)
}

func TestFindImageIdByHashInternalErrShouldFail(t *testing.T) {
	repo, db := BaseSetup()
	db.Close()
	_, err := repo.FindImageIdByHash(nil)
	assert.ErrorIs(t, err, e.ErrInternal)
}
