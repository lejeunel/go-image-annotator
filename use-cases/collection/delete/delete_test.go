package delete

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"testing"
)

func TestDeleteNonExistingCollectionShouldFail(t *testing.T) {

	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Missing: true})
	itr.Execute(Request{Name: name}, p)
	if !p.GotNotFoundErr {
		t.Fatal("expected not found error, but got none")
	}
	if p.GotSuccess {
		t.Fatal("expected no success")
	}
}

func TestDeleteCollectionWithAssociatedResourcesShouldFail(t *testing.T) {

	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{IsPopulated_: true})
	itr.Execute(Request{Name: name}, p)
	if !p.GotDependencyErr {
		t.Fatal("expected dependency error, but got none")
	}
	if p.GotSuccess {
		t.Fatal("expected no success")
	}
}

func TestHandleInternalErrorOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnDelete: true, Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr {
		t.Fatal("expected internal error, but got none")
	}
}

func TestDeleteCollection(t *testing.T) {

	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(Request{Name: name}, p)
	if !p.GotSuccess {
		t.Fatal("expected success, but did not")
	}
}
