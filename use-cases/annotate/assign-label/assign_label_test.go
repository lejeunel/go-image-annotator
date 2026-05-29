package assign_label

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	stest "github.com/lejeunel/go-image-annotator/shared/testing"
)

func TestHandleNotFoundErrOnImageRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected not found error")
	}
}

func TestHandleInternalErrOnImageRetrieval(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestAssignNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{MissingLabel: true}, &st.FakeImageStore{})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected not found error")
	}
}
func TestInternalErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFindLabel: true, Err: e.ErrInternal}, &st.FakeImageStore{})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestAssignLabelToImage(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), *collection)
	label := lbl.NewLabel(lbl.NewLabelId(), "al-label")
	req := Request{ImageId: image.Id, Collection: collection.Name, Label: label.Name}
	repo := &FakeRepo{ReturnLabel: *label}
	itr := NewInteractor(repo, &st.FakeImageStore{Return: image})
	itr.Execute(req, p)
	resp := p.Got
	stest.AssertEqual(t, "label", resp.Label, req.Label)
	stest.AssertEqual(t, "collection", resp.Collection, req.Collection)
	stest.AssertEqual(t, "image id", resp.ImageId, req.ImageId)
	stest.AssertEqual(t, "added label id", repo.AddedLabelId, label.Id)
	stest.AssertEqual(t, "added on image id", repo.AddedOnImageId, image.Id)
	stest.AssertEqual(t, "added on collection id", repo.AddedOnCollectionId, collection.Id)
}
