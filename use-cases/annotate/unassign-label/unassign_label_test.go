package unassign_label

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestHandleNotFoundErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected not found error")
	}
}

func TestHandleInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, &st.FakeImageStore{Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error")
	}
}

func TestUnassignLabelNotAssignedToImageShouldFail(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	image := im.NewImage(im.NewImageId(), clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := NewInteractor(repo, &st.FakeImageStore{Return: &image})
	itr.Execute(Request{Label: "label-not-assigned-to-image"}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatal("expected not found error")
	}
}

func TestHandleInternalErrOnRemoveLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	image := im.NewImage(im.NewImageId(), clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	repo := &FakeRepo{ErrOnRemoveLabel: true, Err: e.ErrInternal}
	itr := NewInteractor(repo, &st.FakeImageStore{Return: &image})
	itr.Execute(Request{Label: "a-label"}, p)
	if p.GotSuccess || !p.GotInternalErr {
		t.Fatal("expected internal error")
	}
}
