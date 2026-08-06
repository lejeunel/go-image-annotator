package update

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
	itr := New(&fk.CollectionRepo{ReturnGroup: group.Name},
		&fk.MetaDataRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}
func TestCheckExistenceOfKeyError(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{ErrOnKeyExists: e.ErrInternal})
	p := &FakePresenter{}
	key, value := "the-key", "the-value"
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: value},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestNonExistingKeyShouldFail(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{})
	p := &FakePresenter{}
	key := "the-key"
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: collection.Name,
			Key: key},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestErrorOnGetCurrentValue(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	itr := New(&fk.CollectionRepo{}, &fk.MetaDataRepo{ErrOnGet: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestErrorOnMismatchingValueTypes(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	key := "the-key"
	itr := New(&fk.CollectionRepo{},
		&fk.MetaDataRepo{ExistingKeys: []string{key}, ReturnValue: "hello"})
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: collection.Name,
			Key: key, Value: 123},
		p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotValidationErr)
}

func TestErrorOnUpdate(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	key, value := "the-key", "the-value"
	newValue := "the-new-value"
	m := &fk.MetaDataRepo{
		ExistingKeys: []string{key},
		ReturnValue:  value, ErrOnUpdate: e.ErrInternal}
	itr := New(&fk.CollectionRepo{}, m)
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: newValue},
		p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestUpdate(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	key, value := "the-key", "the-value"
	newValue := "the-new-value"
	m := &fk.MetaDataRepo{
		ExistingKeys: []string{key},
		ReturnValue:  value}
	itr := New(&fk.CollectionRepo{}, m)
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: newValue},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, m.UpdatedKey, key)
	assert.Equal(t, m.UpdatedValue, newValue)
}
