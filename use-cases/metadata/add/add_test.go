package add

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
	itr := New(&fk.CollectionRepo{}, &fk.ImageRepo{},
		&fk.MetaDataRepo{},
		&fk.StringValidator{},
		&fk.ValueValidator{})

	return itr, collection, image, group
}

func TestHandleAuthError(t *testing.T) {
	itr, collection, image, group := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ReturnGroup: group.Name}
	itr.Auth = &fk.Auth{Err: e.ErrAuthorization}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestKeyValidation(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.KeyValidator = &fk.StringValidator{Invalid: true}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestValueValidation(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.ValueValidator = &fk.ValueValidator{Invalid: true}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestCheckExistenceOfKeyError(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.MetaDataRepo = &fk.MetaDataRepo{ErrOnKeyExists: e.ErrInternal}
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

func TestExistingKeyShouldFail(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.MetaDataRepo = &fk.MetaDataRepo{ExistingKeys: []string{"the-key"}}
	p := &FakePresenter{}
	key, value := "the-key", "the-value"
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: value,
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestErrorOnCheckExistenceOfCollectionShouldFail(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ErrOnExists: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: "key", Value: "value",
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	itr, collection, image, _ := Setup()
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: "key", Value: "value",
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestErrorOnCheckImageIsInCollectionShouldFail(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"my-collection"}}
	itr.ImageRepo = &fk.ImageRepo{ErrOnImageExistsInCollection: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: "key", Value: "value",
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestShouldFailWhenImageIsNotMemberOfCollection(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"my-collection"}}
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: false}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: "key", Value: "value",
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestErrorAddMetaData(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"my-collection"}}
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: true}
	itr.MetaDataRepo = &fk.MetaDataRepo{ErrOnAdd: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotInternalErr)
}

func TestAddMetaData(t *testing.T) {
	itr, collection, image, _ := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{"my-collection"}}
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: true}
	m := &fk.MetaDataRepo{}
	itr.MetaDataRepo = m
	p := &FakePresenter{}
	key, value := "the-key", "the-value"
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key, Value: value,
		},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, m.AddedKey, key)
	assert.Equal(t, m.AddedValue, value)
}
