package read

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func Setup() (Interactor, clc.Collection, im.Image) {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	image := im.NewImage(im.NewImageId(), collection)
	return New(&fk.CollectionRepo{ExistingNames: []string{collection.Name}},
		&fk.ImageRepo{ImageIsInCollection: true},
		&fk.MetaDataRepo{}), collection, image
}

func TestCheckExistenceOfCollectionError(t *testing.T) {
	itr, collection, image := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ErrOnExists: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestMissingCollectionShouldFail(t *testing.T) {
	itr, collection, image := Setup()
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestCheckImageInCollectionError(t *testing.T) {
	itr, collection, image := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	itr.ImageRepo = &fk.ImageRepo{ErrOnImageExistsInCollection: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestImageNotInCollectionShouldFail(t *testing.T) {
	itr, collection, image := Setup()
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: false}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: image.Id.String(), Collection: collection.Name},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrValidation)
}

func TestCheckExistenceOfKeyError(t *testing.T) {
	itr, collection, image := Setup()
	itr.MetaDataRepo = &fk.MetaDataRepo{ErrOnKeyExists: e.ErrInternal}
	itr.ImageRepo = &fk.ImageRepo{ImageIsInCollection: true}
	itr.CollectionRepo = &fk.CollectionRepo{ExistingNames: []string{collection.Name}}
	p := &FakePresenter{}
	key := "the-key"
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key,
		},
		p)
	assert.ErrorIs(t, p.GotErr, e.ErrInternal)
}

func TestNonExistingKeyShouldFail(t *testing.T) {
	itr, collection, image := Setup()
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
	itr, collection, _ := Setup()
	itr.MetaDataRepo = &fk.MetaDataRepo{ErrOnGet: e.ErrInternal}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: collection.Name},
		p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestRead(t *testing.T) {
	itr, collection, image := Setup()
	value := 123
	key := "the-key"
	itr.MetaDataRepo = &fk.MetaDataRepo{ReturnValue: value, ExistingKeys: []string{key}}
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId: image.Id.String(), Collection: collection.Name,
			Key: key,
		},
		p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, value, p.Got.Data.Value)
}
