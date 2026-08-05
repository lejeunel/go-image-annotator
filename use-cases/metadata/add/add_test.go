package add

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleAuthError(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	group := g.NewGroup(g.NewGroupId(), "my-group")
	image.Collection.Group = &group.Name
	itr := New(&fk.CollectionRepo{ReturnGroup: group.Name}, &fk.MetaDataRepo{},
		&fk.StringValidator{},
		&fk.ValueValidator{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestKeyValidation(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{},
		&fk.StringValidator{Invalid: true},
		&fk.ValueValidator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestValueValidation(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{},
		&fk.StringValidator{}, &fk.ValueValidator{Invalid: true})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestErrorAddMetaData(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{ErrOnAdd: e.ErrInternal},
		&fk.StringValidator{}, &fk.ValueValidator{})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestCheckExistenceOfKeyError(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{ErrOnKeyExists: e.ErrInternal},
		&fk.StringValidator{}, &fk.ValueValidator{})
	p := &FakePresenter{}
	key, value := "the-key", "the-value"
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: value},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestExistingKeyShouldFail(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{ExistingKeys: []string{"the-key"}},
		&fk.StringValidator{}, &fk.ValueValidator{})
	p := &FakePresenter{}
	key, value := "the-key", "the-value"
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: value},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestAddMetaData(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	m := &fk.MetaDataRepo{}
	itr := New(&fk.CollectionRepo{}, m,
		&fk.StringValidator{}, &fk.ValueValidator{})
	p := &FakePresenter{}
	key, value := "the-key", "the-value"
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: value},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, m.AddedKey, key)
	assert.Equal(t, m.AddedValue, value)
}
