package image

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	st "github.com/lejeunel/go-image-annotator/shared/testing"
	"github.com/stretchr/testify/assert"
)

func TestIterateImages(t *testing.T) {
	imr, cr, _ := SetupList()
	collectionName := "a-collection"
	collection := clc.NewCollection(clc.NewCollectionId(), collectionName)
	cr.Create(collection)
	im0 := CreateImageInCollectionFromString(
		imr,
		collection,
		st.FakeUUIDFromInt(0),
	)
	im1 := CreateImageInCollectionFromString(
		imr,
		collection,
		st.FakeUUIDFromInt(1),
	)

	res := []im.BaseImage{}
	for batch, err := range imr.Iterate("", 1) {
		assert.NoError(t, err)
		res = append(res, batch)
	}
	assert.Equal(t, 2, len(res))
	assert.Equal(t, im0.Id.String(), res[0].ImageId.String())
	assert.Equal(t, im1.Id.String(), res[1].ImageId.String())
}
