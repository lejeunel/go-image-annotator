package assign_label

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleNotFoundErrOnImageRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrNotFound})
	itr.Execute(Request{im.NewImageId().String(), "a-collection", "a-label"}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected not found error")
	}
}

func TestHandleInternalErrOnImageRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrInternal})
	itr.Execute(Request{im.NewImageId().String(), "a-collection", "a-label"}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestAssignNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := NewInteractor(&FakeRepo{MissingLabel: true}, &st.FakeImageStore{Return: &image})
	itr.Execute(Request{image.Id.String(), collection.Name, "a-label"}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected not found error")
	}
}
func TestInternalErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	itr := NewInteractor(&FakeRepo{ErrOnFindLabel: true, Err: e.ErrInternal}, &st.FakeImageStore{Return: &image})
	itr.Execute(Request{im.NewImageId().String(), "a-collection", "a-label"}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestAssignLabelToImage(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	label := lbl.NewLabel(lbl.NewLabelId(), "al-label")
	req := Request{ImageId: image.Id.String(), Collection: collection.Name, Label: label.Name}
	repo := &FakeRepo{ReturnLabel: label}
	itr := NewInteractor(repo, &st.FakeImageStore{Return: &image})
	itr.Execute(req, p)
	resp := p.Got
	assert.Equal(t, resp.Label, req.Label, "label")
	assert.Equal(t, resp.Collection, req.Collection, "collection")
	assert.Equal(t, resp.ImageId, req.ImageId, "image id")
	assert.Equal(t, repo.AddedLabelId, label.Id, "added label id")
	assert.Equal(t, repo.AddedOnImageId, image.Id, "added on image id")
	assert.Equal(t, repo.AddedOnCollectionId, collection.Id, "added on collection id")
}
