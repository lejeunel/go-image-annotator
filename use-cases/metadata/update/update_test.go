package update

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Setup() (Interactor, clc.Collection, im.Image, g.Group) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	group := g.NewGroup(g.NewGroupId(), "my-group")
	image.Collection.Group = &group.Name
	return New(&fk.CollectionRepo{}, &fk.ImageRepo{},
		&fk.MetaDataRepo{}), collection, image, group
}

func TestHandleAuthError(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.Auth = &fk.Auth{Err: e.ErrAuthorization}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCheckExistenceOfCollectionError(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ErrOnExists: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestMissingCollectionShouldFail(t *testing.T) {
	itr, collection, image, _ := Setup()
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestCheckImageInCollectionError(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	itr.ImageRepo = &fk.ImageRepo{ErrOnImageExistsInCollection: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestImageNotInCollectionShouldFail(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: false}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestCheckExistenceOfKeyError(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.MetaDataRepo = &fk.MetaDataRepo{ErrOnKeyExists: e.ErrInternal}
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: true}
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	p := &FakePresenter{}
	key, value := "the-key", "the-value"
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: value,
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestNonExistingKeyShouldFail(t *testing.T) {
	itr, collection, image, _ := Setup()
	p := &FakePresenter{}
	key := "the-key"
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key,
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestErrorOnGetCurrentValue(t *testing.T) {
	itr, collection, _, _ := Setup()
	itr.MetaDataRepo = &fk.MetaDataRepo{ErrOnGet: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestErrorOnMismatchingValueTypes(t *testing.T) {
	itr, collection, _, _ := Setup()
	key := "the-key"
	itr.MetaDataRepo = &fk.MetaDataRepo{ExistingKeys: []string{key}, ReturnValue: "hello"}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: im.NewImageId().String(), Collection: collection.Name,
			Key: key, Value: 123,
		},
		p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotValidationErr)
}

func TestErrorOnUpdate(t *testing.T) {
	itr, collection, image, _ := Setup()
	key, value := "the-key", "the-value"
	newValue := "the-new-value"
	m := &fk.MetaDataRepo{
		ExistingKeys: []string{key},
		ReturnValue:  value, ErrOnUpdate: e.ErrInternal,
	}
	itr.MetaDataRepo = m
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: true}
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: newValue,
		},
		p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestUpdate(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: true}
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	key, value := "the-key", "the-value"
	newValue := "the-new-value"
	m := &fk.MetaDataRepo{
		ExistingKeys: []string{key},
		ReturnValue:  value,
	}
	itr.MetaDataRepo = m
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: newValue,
		},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, m.UpdatedKey, key)
	assert.Equal(t, m.UpdatedValue, newValue)
}
