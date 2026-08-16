package image

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/testing"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnImageIsInCollectionShouldFail(t *testing.T) {
	imRepo, _, db := SetupAdd(s.NewInMemory())
	db.Close()
	_, err := imRepo.ImageExistsInCollection(im.NewImageId(), "a-collection")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddedImageToCollectionExists(t *testing.T) {
	imRepo, clcRepo, _ := SetupAdd(s.NewInMemory())
	imageId, _, _ := AddToCollection(imRepo, clcRepo, "a-collection", "the-hash")
	isAdded, err := imRepo.ImageExistsInCollection(*imageId, "a-collection")
	assert.NoError(t, err)
	assert.True(t, isAdded)
}

func TestInternalErrOnImageExistsShouldFail(t *testing.T) {
	db := s.NewInMemory()
	imr, _, db := SetupAdd(db)
	db.Close()
	_, err := imr.ImageExists(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}
