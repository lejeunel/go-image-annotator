package update_label

import (
	"errors"
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	stest "github.com/lejeunel/go-image-annotator/shared/testing"
)

func TestHandleErrOnFindLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrNotFound, ErrOnFindLabel: true})
	itr.Execute(Request{}, p)
	if !errors.Is(p.GotErr, e.ErrNotFound) || p.GotSuccess {
		t.Fatalf("expected to get error %v, got %v, with success %v",
			e.ErrNotFound, p.GotErr, p.GotSuccess)
	}
}

func TestHandleErrOnUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	itr := NewInteractor(&FakeRepo{Returns: label, Err: e.ErrNotFound, ErrOnUpdate: true})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get error %v, got %v, with success %v",
			e.ErrNotFound, p.GotErr, p.GotSuccess)
	}
}

func TestFetchLabelFromName(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: newLabel}
	itr := NewInteractor(repo)
	itr.Execute(Request{AnnotationId: a.NewAnnotationId(), Label: newLabel.Name}, p)
	stest.AssertEqual(t, "label name", repo.FetchedLabelWithName, newLabel.Name)
}

func TestUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: newLabel}
	itr := NewInteractor(repo)
	annotationId := a.NewAnnotationId()
	itr.Execute(Request{AnnotationId: annotationId, Label: newLabel.Name}, p)
	stest.AssertEqual(t, "annotation id", repo.UpdatedAnnotationId, annotationId)
	stest.AssertEqual(t, "label id", repo.UpdatedLabelId, newLabel.Id)
}
