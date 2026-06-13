package delete

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/modules/image-store"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := g.NewGroup(g.NewGroupId(), "my-group")
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection",
		clc.WithGroup(group))
	image := im.NewImage(im.NewImageId(), collection)
	itr := NewInteractor(&st.FakeImageStore{Return: &image}, &FakeRepo{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{ImageId: im.NewImageId().String(), Collection: "a-collection"},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingResourceShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrNotFound}, &FakeRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrInternal}, &FakeRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrNotFound})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnRemoveLabel(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteNonExistingBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1,
		lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	image.AddBoundingBox(box)
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrNotFound})
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
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnRemoveImageFromCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveImage: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestRemoveImageFromCollection(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotSuccess)
}
