package image

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleInternalErrOnDeleteImage(t *testing.T) {
	imr, db := BaseSetup()
	db.Close()
	err := imr.Delete(im.NewImageId())
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestDeleteImage(t *testing.T) {
	imr, _ := BaseSetup()
	id := im.NewImageId()
	imr.AddImage(id, nil, im.Specs{})
	err := imr.Delete(id)
	assert.NoError(t, err)
}
