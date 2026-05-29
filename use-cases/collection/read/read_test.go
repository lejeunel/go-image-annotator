package read

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestReadCollection(t *testing.T) {
	name := "my-collection"
	desc := "a-description"
	repo := &FakeRepo{Collection: clc.Collection{Name: name, Description: desc}}
	p := &FakeReadPresenter{}
	itr := NewInteractor(repo)
	req := Request{Name: name}
	want := Response{Name: name, Description: desc}
	itr.Execute(req, p)
	if p.Got != want {
		t.Fatalf("expected %v, got %v", want, p.Got)
	}
}

func TestReadNonExistingCollectionShouldFail(t *testing.T) {
	repo := &FakeRepo{Collection: clc.Collection{Name: "my-collection", Description: "a-description"}}
	p := &FakeReadPresenter{}
	itr := NewInteractor(repo)
	req := Request{Name: "non-existing-collection"}
	itr.Execute(req, p)
	if !p.GotNotFoundErr {
		t.Fatal("expected not found error, but got none")
	}
	if p.GotSuccess {
		t.Fatal("expected no success")
	}
}

func TestHandleInternalError(t *testing.T) {
	p := &FakeReadPresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr {
		t.Fatal("expected internal error, but got none")
	}
}
