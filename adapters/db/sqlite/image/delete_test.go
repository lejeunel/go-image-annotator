package image

import (
	"testing"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnDeleteImage(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteImageRepo(db)
	db.Close()
	err := repo.Delete(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteImage(t *testing.T) {
	repo := NewSQLiteImageRepo(s.NewInMemory())
	id := im.NewImageId()
	repo.AddImage(id, nil, im.ImageSpecs{})
	err := repo.Delete(id)
	assert.NoError(t, err)
}
