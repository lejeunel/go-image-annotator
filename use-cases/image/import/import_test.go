package import_image

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection",
		clc.WithGroup(g.NewGroup(g.NewGroupId(), "dst-group")))
	itr := NewInteractor(&FakeRepo{Return: dstCollection},
		WithAuth(FailingAuth{}))
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
	itr := NewInteractor(&FakeRepo{ImageMissing: true})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnFindingSourceImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnImageExists: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingDestinationCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFindCollection: true, Err: e.ErrNotFound})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnFindCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFindCollection: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}

	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection",
		clc.WithGroup(g.NewGroup(g.NewGroupId(), "dst-group")))
	itr := NewInteractor(&FakeRepo{Return: dstCollection, ImageAlreadyInCollection: true})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnImageAlreadyExistsInCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	itr := NewInteractor(&FakeRepo{Return: dstCollection, ErrOnImageExistsInCollection: true,
		Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
func TestInternalErrOnImportShouldFail(t *testing.T) {
	p := &FakePresenter{}
	dstCollection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	repo := &FakeRepo{Return: dstCollection, ErrOnImport: true, Err: e.ErrInternal}
	itr := NewInteractor(repo)
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String()}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestImportImageInCollection(t *testing.T) {
	p := &FakePresenter{}
	imageId := im.NewImageId()
	collection := clc.NewCollection(clc.NewCollectionId(), "dst-collection")
	repo := &FakeRepo{Return: collection}
	itr := NewInteractor(repo)
	itr.Execute(t.Context(),
		Request{ImageId: imageId.String(),
			SourceCollection:      "src-collection",
			DestinationCollection: collection.Name}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, imageId, repo.ImportedImageId)
	assert.Equal(t, collection.Id, repo.ImportedIntoCollectionId)
}
