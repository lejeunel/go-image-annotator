package delete

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestDeleteLabelWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{IsUsed_: true})
	itr.Execute(Request{}, p)
	if !p.GotDependencyErr {
		t.Fatal("expected dependency error, but got none")
	}
	if p.GotSuccess {
		t.Fatal("expected no success")
	}
}

func TestHandleInternalErrOnIsUsed(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal, ErrOnIsUsed: true})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !p.GotInternalErr {
		t.Fatal("expected internal error")
	}
}

func TestHandleInternalErrOnExists(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal, ErrOnExists: true})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !p.GotInternalErr {
		t.Fatal("expected internal error")
	}
}

func TestDeletingMissingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{IsMissing: true})
	itr.Execute(Request{}, p)
	if p.GotSuccess || !p.GotNotFoundErr {
		t.Fatal("expected not found error")
	}
}

func TestDeleteLabel(t *testing.T) {

	name := "my-collection"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(Request{Name: name}, p)
	if !p.GotSuccess {
		t.Fatal("expected success, but did not")
	}
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr {
		t.Fatal("expected internal error, but got none")
	}
}
