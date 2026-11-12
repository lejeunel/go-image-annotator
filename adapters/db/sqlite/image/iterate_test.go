package image

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
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

func addImageFromString(t *testing.T, imr ImageRepo, id string) im.ImageId {
	imageId, err := im.NewImageIdFromString(id)
	assert.NoError(t, err)
	assert.NoError(t, imr.AddImage(imageId, []byte(id), im.Specs{}))
	return imageId
}

func collectIDs(t *testing.T, imr ImageRepo, pageSize int) []string {
	ids := []string{}
	for id, err := range imr.IterateIDs(pageSize) {
		assert.NoError(t, err)
		ids = append(ids, id.String())
	}
	return ids
}

func TestIterateIDsYieldsEachImageOnce(t *testing.T) {
	imr, cr, _ := SetupList()
	a := clc.NewCollection(clc.NewCollectionId(), "a")
	b := clc.NewCollection(clc.NewCollectionId(), "b")
	cr.Create(a)
	cr.Create(b)
	shared := addImageFromString(t, imr, st.FakeUUIDFromInt(0))
	assert.NoError(t, imr.AddToCollection(shared, a.Name))
	assert.NoError(t, imr.AddToCollection(shared, b.Name))
	inA := addImageFromString(t, imr, st.FakeUUIDFromInt(1))
	assert.NoError(t, imr.AddToCollection(inA, a.Name))
	orphan := addImageFromString(t, imr, st.FakeUUIDFromInt(2))

	want := []string{shared.String(), inA.String(), orphan.String()}
	for _, pageSize := range []int{1, 2, 3, 10} {
		assert.Equal(t, want, collectIDs(t, imr, pageSize), "page size %v", pageSize)
	}
}

func TestIterateIDsEmpty(t *testing.T) {
	imr, _, _ := SetupList()
	assert.Empty(t, collectIDs(t, imr, 2))
}

func TestIterateIDsStopsEarly(t *testing.T) {
	imr, _, _ := SetupList()
	for i := range 3 {
		addImageFromString(t, imr, st.FakeUUIDFromInt(i))
	}
	n := 0
	for _, err := range imr.IterateIDs(1) {
		assert.NoError(t, err)
		n++
		if n == 2 {
			break
		}
	}
	assert.Equal(t, 2, n)
}

func TestIterateIDsInvalidPageSize(t *testing.T) {
	imr, _, _ := SetupList()
	for _, err := range imr.IterateIDs(0) {
		assert.ErrorIs(t, err, e.ErrValidation)
	}
}
