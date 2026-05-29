package delete

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestNonExistingResourceShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrNotFound}, &FakeRepo{})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestHandleInternalErr(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrInternal}, &FakeRepo{})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected to get internal error")
	}
}

func TestDeleteNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !(p.GotNotFoundErr) {
		t.Fatalf("expected not found error")
	}
}

func TestHandleInternalErrOnRemoveLabel(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	image.AddLabel(lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !(p.GotInternalErr) {
		t.Fatalf("expected internal error")
	}
}

func TestDeleteNonExistingBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1,
		lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	image.AddBoundingBox(*box)
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !(p.GotNotFoundErr) {
		t.Fatalf("expected not found error")
	}
}

func TestHandleInternalErrOnDeleteBoxes(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	box := a.NewBoundingBox(a.NewAnnotationId(), 1, 1, 1, 1, lbl.NewLabel(lbl.NewLabelId(), "a-label"))
	image.AddBoundingBox(*box)
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveAnnotation: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !(p.GotInternalErr) {
		t.Fatalf("expected internal error")
	}
}

func TestInternalErrOnRemoveImageFromCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{ErrOnRemoveImage: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !(p.GotInternalErr) {
		t.Fatalf("expected internal error")
	}
}

func TestRemoveImageFromCollection(t *testing.T) {
	p := &FakePresenter{}
	id := im.NewImageId()
	image := im.NewImage(id, clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := NewInteractor(&st.FakeImageStore{Return: &image},
		&FakeRepo{})
	itr.Execute(Request{}, p)
	if !p.GotSuccess {
		t.Fatalf("expected success")
	}
}
