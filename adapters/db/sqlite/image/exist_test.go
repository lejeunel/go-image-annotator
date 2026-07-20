package image

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalErrOnImageIsInCollectionShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteImageRepo(db)
	db.Close()
	_, err := repo.ImageExistsInCollection(im.NewImageId(), clc.NewCollectionId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestAddedImageToCollectionExists(t *testing.T) {
	imRepo, clcRepo := MakeRepos(s.NewInMemory())
	imageId, collectionId, _ := AddToCollection(imRepo, clcRepo, "a-collection", "the-hash")
	isAdded, err := imRepo.ImageExistsInCollection(*imageId, *collectionId)
	assert.NoError(t, err)
	assert.True(t, isAdded)
}

func TestInternalErrOnImageExistsShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteImageRepo(db)
	db.Close()
	_, err := repo.ImageExists(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}
