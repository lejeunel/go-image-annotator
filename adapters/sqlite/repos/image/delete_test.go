package image

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnDeleteImage(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	repo.Db.Close()
	err := repo.Delete(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteImage(t *testing.T) {
	repo := NewTestSQLiteImageRepo()
	id := im.NewImageId()
	repo.AddImage(id, nil, im.ImageSpecs{})
	err := repo.Delete(id)
	assert.NoError(t, err)
}
