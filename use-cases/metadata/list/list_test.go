package list

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorOnListShouldFail(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.MetaDataRepo{ErrOnList: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestList(t *testing.T) {
	items := []m.MetaData{
		{Key: "first-key", Value: "hello"},
		{Key: "second-key", Value: 123},
	}
	itr := New(&fk.MetaDataRepo{ReturnList: items})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: "my-collection"},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, items, p.Got.MetaData)
}
