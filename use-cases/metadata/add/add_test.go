package add

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Setup() (Interactor, im.Image) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	group := g.NewGroup(g.NewGroupId(), "my-group")
	image.Collection.Group = &group.Name
	itr := New(&fk.ImageStore{},
		&fk.MetaDataRepo{},
		&fk.StringValidator{},
		&fk.ValueValidator{})

	return itr, image
}

func TestHandleAuthError(t *testing.T) {
	itr, image := Setup()
	itr.ImageStore = &fk.ImageStore{Return: &image}
	itr.Auth = &fk.Auth{ErrOnAuth: e.ErrAuthorization}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: image.Collection.Name},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestKeyValidation(t *testing.T) {
	itr, image := Setup()
	itr.KeyValidator = &fk.StringValidator{Invalid: true}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: image.Collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestValueValidation(t *testing.T) {
	itr, image := Setup()
	itr.ValueValidator = &fk.ValueValidator{Invalid: true}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: image.Collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestErrorAddMetaData(t *testing.T) {
	itr, image := Setup()
	itr.MetaDataRepo = &fk.MetaDataRepo{ErrOnAdd: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: image.Collection.Name},
		p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestAddMetaData(t *testing.T) {
	itr, image := Setup()
	key, value := "the-key", "the-value"
	m := &fk.MetaDataRepo{ReturnList: []m.MetaData{{Key: key, Value: value}}}
	itr.MetaDataRepo = m
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: image.Collection.Name,
			Key: key, Value: value,
		},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, m.AddedKey, key)
	assert.Equal(t, m.AddedValue, value)
	assert.Equal(t, 1, len(p.Got.Meta))
}
