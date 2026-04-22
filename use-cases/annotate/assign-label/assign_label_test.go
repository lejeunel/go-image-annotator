package assign_label

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
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
	if !p.GotSuccess || !(resp.Label == req.Label) || !(resp.Collection == req.Collection) || !(resp.ImageId == req.ImageId) {
		t.Fatalf("expected response does not match, got request %+v and response %+v", req, resp)
	}
	if (repo.AddedLabelId != label.Id) || (repo.AddedOnImageId != image.Id) || (repo.AddedOnCollectionId != collection.Id) {
		t.Fatalf("expected to add label %v on image %v and collection %v, got %v, %v, %v",
			label, image.Id, image.Collection.Id,
			repo.AddedLabelId, repo.AddedOnImageId, repo.AddedOnCollectionId)
	}
}
