package read

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrorOnFind(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&st.FakeImageStore{Err: e.ErrNotFound})
	itr.Execute(Request{ImageId: im.NewImageId().String(), Collection: "a-collection"}, p)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
	assert.False(t, p.GotSuccess)
}

func TestFindImageGivesCorrectIdAndCollection(t *testing.T) {
	p := &FakePresenter{}
	existingImage := im.NewImage(im.NewImageId(), clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := New(&st.FakeImageStore{Return: &existingImage})
	itr.Execute(Request{ImageId: existingImage.Id.String(),
		Collection: existingImage.Collection.Name}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, p.Got.Id, existingImage.Id, "id")
	assert.Equal(t, p.Got.Collection.Name, existingImage.Collection.Name, "collection name")
}
