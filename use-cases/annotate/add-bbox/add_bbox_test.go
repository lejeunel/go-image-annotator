package add_bbox

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"github.com/stretchr/testify/assert"
)

func CreateImage() im.Image {
	collection := clc.NewCollection(clc.NewCollectionId(), "my-collection")
	return im.NewImage(im.NewImageId(), collection)
}

func TestHandleAuthError(t *testing.T) {
	image := CreateImage()
	group := g.NewGroup(g.NewGroupId(), "my-group")
	image.Collection.Group = &group
	itr := NewInteractor(&st.FakeImageStore{Return: &image}, &FakeRepo{},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingImageStoreResourceShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrNotFound}, &FakeRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnImageRetrievalShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrInternal}, &FakeRepo{})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestNotFoundErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnFindLabel: true, Err: e.ErrNotFound})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnFindLabel: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestValidationErrOnAddBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String(), Collection: "a-collection", Label: "a-label",
		Xc: 1, Yc: 1, Width: -999, Height: 3}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestNotFoundErrOnAddBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnAdd: true, Err: e.ErrNotFound})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String(), Collection: "a-collection", Label: "a-label",
		Xc: 1, Yc: 1, Width: 3, Height: 3}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnAddBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnAdd: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{ImageId: im.NewImageId().String(), Collection: "a-collection", Label: "a-label",
		Xc: 1, Yc: 1, Width: 3, Height: 3}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAddBoundingBox(t *testing.T) {
	p := &FakePresenter{}
	repo := FakeRepo{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	req := Request{ImageId: image.Id.String(), Collection: collection.Name,
		Label: "a-label", Xc: float32(1.0), Yc: float32(1.0), Width: float32(3.0),
		Height: float32(3.0)}
	itr := NewInteractor(&st.FakeImageStore{Return: &image}, &repo)
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, repo.GotImageId.String(), req.ImageId, "image id")
	assert.Equal(t, repo.GotCollectionId, collection.Id, "collection id")
	assert.Equal(t, repo.GotBox.Label.Name, req.Label, "label name")
	assert.Equal(t, repo.GotBox.Xc, req.Xc, "xc")
	assert.Equal(t, repo.GotBox.Yc, req.Yc, "yc")
	assert.Equal(t, repo.GotBox.Width, req.Width, "width")
	assert.Equal(t, repo.GotBox.Height, req.Height, "height")

}
