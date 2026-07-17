package delete

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "my-group")
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection",
		clc.WithGroup(group))
	image := im.NewImage(im.NewImageId(), collection)
	itr := New(&fk.ImageStore{Return: &image}, &fk.ImageRepo{},
		&fk.AnnotationRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: "a-collection"},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingResourceShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageStore{Err: e.ErrNotFound}, &fk.ImageRepo{}, &fk.AnnotationRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ImageStore{Err: e.ErrInternal}, &fk.ImageRepo{}, &fk.AnnotationRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	itr := New(&fk.ImageStore{Return: &image},
		&fk.ImageRepo{}, &fk.AnnotationRepo{ErrOnRemoveAnnotation: e.ErrNotFound})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1,
		lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	image.AddBoundingBox(box)
	itr := New(&fk.ImageStore{Return: &image},
		&fk.ImageRepo{}, &fk.AnnotationRepo{ErrOnRemoveAnnotation: e.ErrNotFound})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnDeleteBoxes(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	image.AddBoundingBox(box)
	itr := New(&fk.ImageStore{Return: &image},
		&fk.ImageRepo{}, &fk.AnnotationRepo{ErrOnRemoveAnnotation: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnRemoveImageFromCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := New(&fk.ImageStore{Return: &image},
		&fk.ImageRepo{ErrOnRemoveImage: e.ErrInternal}, &fk.AnnotationRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestRemoveImageFromCollection(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := New(&fk.ImageStore{Return: &image},
		&fk.ImageRepo{}, &fk.AnnotationRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotSuccess)
}
