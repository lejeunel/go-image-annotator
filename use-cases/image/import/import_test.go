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
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection",
		clc.WithGroup(g.NewGroup(g.NewGroupId(), "dst-group")))
	itr := New(&fk.ImageRepo{}, &fk.CollectionRepo{Return: dstCollection},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(),
			SourceCollection:      "src-collection",
			DestinationCollection: "dst-collection"},
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

	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection",
		clc.WithGroup(g.NewGroup(g.NewGroupId(), "dst-group")))
	itr := New(&fk.ImageRepo{ImageAlreadyInCollection: true}, &fk.CollectionRepo{Return: dstCollection})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	itr := New(&fk.ImageRepo{ErrOnImageExistsInCollection: e.ErrInternal}, &fk.CollectionRepo{Return: dstCollection})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
func TestInternalErrOnImportShouldFail(t *testing.T) {
	p := &FakePresenter{}
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	itr := New(&fk.ImageRepo{ErrOnAddToCollection: e.ErrInternal}, &fk.CollectionRepo{Return: dstCollection})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestImportImageInCollection(t *testing.T) {
	p := &FakePresenter{}
	imageId := im.NewImageId()
	collection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	repo := &fk.ImageRepo{}
	itr := New(repo, &fk.CollectionRepo{Return: collection})
	itr.Execute(t.Context(),
		Request{ImageId: imageId.String(),
			SourceCollection:      "src-collection",
			DestinationCollection: collection.Name}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, imageId, repo.ImportedImageId)
	assert.Equal(t, collection.Id, repo.ImportedIntoCollectionId)
}
