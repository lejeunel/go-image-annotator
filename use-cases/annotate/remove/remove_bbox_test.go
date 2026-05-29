package remove

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestNonExistingBoxShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrNotFound})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestInternalErrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestRemoveBox(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := NewInteractor(repo)
	annotationId := a.NewAnnotationId()
	itr.Execute(Request{Id: annotationId}, p)
	if !p.GotSuccess || !(repo.Got == annotationId) {
		t.Fatalf("expected to remove annotation, got success %v and annotation id %v, wanted %v",
			p.GotSuccess, repo.Got.String(), annotationId.String())
	}
}
