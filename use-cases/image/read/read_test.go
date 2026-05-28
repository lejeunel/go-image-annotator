package read

import (
	"errors"
	"testing"

	st "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	stest "github.com/lejeunel/go-image-annotator-v2/shared/testing"
)

func TestHandleErrorOnImageIdParsing(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrNotFound})
	itr.Execute(Request{ImageId: "invalid-image-id", Collection: "a-collection"}, p)
	if !errors.Is(p.GotErr, e.ErrValidation) || p.GotSuccess {
		t.Fatalf("expected to get validation error")
	}
}

func TestHandleErrorOnFind(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&st.FakeImageStore{Err: e.ErrNotFound})
	itr.Execute(Request{ImageId: im.NewImageId().String(), Collection: "a-collection"}, p)
	if !errors.Is(p.GotErr, e.ErrNotFound) || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestFindImageGivesCorrectIdAndCollection(t *testing.T) {
	p := &FakePresenter{}
	existingImage := im.NewImage(im.NewImageId(), *clc.NewCollection(clc.NewCollectionId(), "a-collection"))
	itr := NewInteractor(&st.FakeImageStore{Return: existingImage})
	itr.Execute(Request{ImageId: existingImage.Id.String(),
		Collection: existingImage.Collection.Name}, p)
	if !p.GotSuccess {
		t.Fatalf("expected to get success")
	}
	stest.AssertEqual(t, "id", p.Got.Id, existingImage.Id)
	stest.AssertEqual(t, "collection name", p.Got.Collection.Name, existingImage.Collection.Name)
}
