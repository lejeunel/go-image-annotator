package add_bbox

import (
	"testing"

	st "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
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
	imageId := im.NewImageId()
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	image := im.NewImage(imageId, *collection)
	label := "a-label"
	x := float32(1.0)
	y := float32(1.0)
	width := float32(3.0)
	height := float32(3.0)
	itr := NewInteractor(&st.FakeImageStore{Return: image}, &repo)
	itr.Execute(Request{ImageId: imageId, Collection: "a-collection", Label: label,
		Xc: x, Yc: y, Width: width, Height: height}, p)
	if !p.GotSuccess {
		t.Fatalf("expected success")
	}
	if repo.GotImageId != imageId {
		t.Fatalf("expected to store box on image %v, got %v", imageId.String(), repo.GotImageId.String())
	}
	if repo.GotCollectionId != collection.Id {
		t.Fatalf("expected to store box on collection %v, got %v", collection.Id.String(), repo.GotCollectionId.String())
	}
	if repo.GotBox.Label.Name != label {
		t.Fatalf("expected to store box with label %v, got %v", label, repo.GotBox.Label.Name)
	}
	if repo.GotBox.Xc != x {
		t.Fatalf("expected to store box with x %v, got %v", x, repo.GotBox.Xc)
	}
	if repo.GotBox.Yc != y {
		t.Fatalf("expected to store box with y %v, got %v", x, repo.GotBox.Yc)
	}
	if repo.GotBox.Width != width {
		t.Fatalf("expected to store box with width %v, got %v", width, repo.GotBox.Width)
	}
	if repo.GotBox.Height != height {
		t.Fatalf("expected to store box with height %v, got %v", height, repo.GotBox.Height)
	}

}
