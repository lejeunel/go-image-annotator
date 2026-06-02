package import_image

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestNonExistingSourceImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ImageMissing: true})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnFindingSourceImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnImageExists: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingDestinationCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFindCollection: true, Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnFindCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFindCollection: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ImageAlreadyInCollection: true})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnImageExistsInCollection: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
func TestInternalErrOnImportShouldFail(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{ErrOnImport: true, Err: e.ErrInternal}
	itr := NewInteractor(repo)
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestImportImageInCollection(t *testing.T) {
	p := &FakePresenter{}
	imageId := im.NewImageId()
	collection := clc.NewCollection(clc.NewCollectionId(), "a-destination-collection")
	repo := &FakeRepo{DestinationCollection: collection}
	itr := NewInteractor(repo)
	itr.Execute(Request{ImageId: imageId, Collection: collection.Name}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, imageId, repo.ImportedImageId)
	assert.Equal(t, collection.Id, repo.ImportedIntoCollectionId)
}
