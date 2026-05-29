package add_bbox

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	stest "github.com/lejeunel/go-image-annotator/shared/testing"
)

func TestNonExistingImageStoreResourceShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrNotFound}, &FakeRepo{})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestInternalErrOnImageRetrievalShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrInternal}, &FakeRepo{})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestNotFoundErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnFindLabel: true, Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected not found error")
	}
}

func TestInternalErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnFindLabel: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestValidationErrOnAddBoxShouldFail(t *testing.T) {
	presenter := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{})
	itr.Execute(Request{ImageId: im.NewImageId(), Collection: "a-collection", Label: "a-label",
		Xc: 1, Yc: 1, Width: -999, Height: 3}, presenter)
	if !presenter.GotValidationErr || presenter.GotSuccess {
		t.Fatalf("expected validation error")
	}
}

func TestNotFoundErrOnAddBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnAdd: true, Err: e.ErrNotFound})
	itr.Execute(Request{ImageId: im.NewImageId(), Collection: "a-collection", Label: "a-label",
		Xc: 1, Yc: 1, Width: 3, Height: 3}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected not found error")
	}
}

func TestInternalErrOnAddBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{}, &FakeRepo{ErrOnAdd: true, Err: e.ErrInternal})
	itr.Execute(Request{ImageId: im.NewImageId(), Collection: "a-collection", Label: "a-label",
		Xc: 1, Yc: 1, Width: 3, Height: 3}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestAddBoundingBox(t *testing.T) {
	p := &FakePresenter{}
	repo := FakeRepo{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(im.NewImageId(), collection)
	req := Request{ImageId: image.Id, Collection: collection.Name,
		Label: "a-label", Xc: float32(1.0), Yc: float32(1.0), Width: float32(3.0),
		Height: float32(3.0)}
	itr := NewInteractor(&st.FakeImageStore{Return: &image}, &repo)
	itr.Execute(req, p)
	if !p.GotSuccess {
		t.Fatalf("expected success")
	}
	stest.AssertEqual(t, "image id", repo.GotImageId, req.ImageId)
	stest.AssertEqual(t, "collection id", repo.GotCollectionId, collection.Id)
	stest.AssertEqual(t, "label name", repo.GotBox.Label.Name, req.Label)
	stest.AssertEqual(t, "xc", repo.GotBox.Xc, req.Xc)
	stest.AssertEqual(t, "yc", repo.GotBox.Yc, req.Yc)
	stest.AssertEqual(t, "width", repo.GotBox.Width, req.Width)
	stest.AssertEqual(t, "height", repo.GotBox.Height, req.Height)

}
