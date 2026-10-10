package import_image

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "dst-group")
	srcCollection := clc.NewCollection(clc.NewCollectionId(), "src-collection",
		clc.WithGroup(group.Name))
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection",
		clc.WithGroup(group.Name))

	itr := New(
		&fk.ImageRepo{},
		&fk.CollectionRepo{Existing: []clc.Collection{srcCollection, dstCollection}},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}),
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{
			ImageId:               im.NewImageId().String(),
			SourceCollection:      srcCollection.Name,
			DestinationCollection: dstCollection.Name,
		},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingSourceImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{ErrOnImageExists: e.ErrNotFound}, &fk.CollectionRepo{})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingDestinationCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageRepo{}, &fk.CollectionRepo{ErrOnFind: e.ErrNotFound})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}

	group := g.NewGroup(g.NewGroupId(), "dst-group")
	source := clc.NewCollection(clc.NewCollectionId(), "source-collection",
		clc.WithGroup(group.Name))
	destination := clc.NewCollection(clc.NewCollectionId(), "destination-collection",
		clc.WithGroup(group.Name))
	itr := New(
		&fk.ImageRepo{ImageIsInCollection: true},
		&fk.CollectionRepo{Existing: []clc.Collection{source, destination}},
	)
	itr.Execute(
		t.Context(),
		Request{
			SourceCollection:      source.Name,
			DestinationCollection: destination.Name,
			ImageId:               im.NewImageId().String(),
		},
		p,
	)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	destination := clc.NewCollection(clc.NewCollectionId(), "destination-collection")
	source := clc.NewCollection(clc.NewCollectionId(), "source-collection")
	itr := New(
		&fk.ImageRepo{ErrOnImageExistsInCollection: e.ErrInternal},
		&fk.CollectionRepo{Existing: []clc.Collection{destination, source}},
	)
	itr.Execute(t.Context(), Request{
		SourceCollection: source.Name, ImageId: im.NewImageId().String(),
		DestinationCollection: destination.Name,
	}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnImportShouldFail(t *testing.T) {
	p := &FakePresenter{}
	destination := clc.NewCollection(clc.NewCollectionId(), "destination-collection")
	source := clc.NewCollection(clc.NewCollectionId(), "source-collection")
	itr := New(
		&fk.ImageRepo{ErrOnAddToCollection: e.ErrInternal},
		&fk.CollectionRepo{Existing: []clc.Collection{source, destination}},
	)
	itr.Execute(t.Context(), Request{
		SourceCollection: source.Name,
		ImageId:          im.NewImageId().String(), DestinationCollection: destination.Name,
	}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestImportImageInCollection(t *testing.T) {
	p := &FakePresenter{}
	imageId := im.NewImageId()
	collection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.CollectionRepo{Existing: []clc.Collection{collection}})
	itr.Execute(t.Context(),
		Request{
			ImageId:               imageId.String(),
			SourceCollection:      "src-collection",
			DestinationCollection: collection.Name,
		}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, imageId, repo.AddedImageId)
	assert.Equal(t, collection.Name, *repo.AddedIntoCollection)
}
