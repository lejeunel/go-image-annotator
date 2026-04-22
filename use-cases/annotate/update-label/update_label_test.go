package update_label

import (
	"errors"
	"testing"

	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator-v2/entities/label"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
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
	if !errors.Is(p.GotErr, e.ErrNotFound) || p.GotSuccess {
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
	if repo.FetchedLabelWithName != newLabel.Name {
		t.Fatalf("expected to fetch label with name %v, got %v",
			newLabel.Name, repo.FetchedLabelWithName)
	}
}

func TestUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: newLabel}
	itr := NewInteractor(repo)
	annotationId := a.NewAnnotationId()
	itr.Execute(Request{AnnotationId: annotationId, Label: newLabel.Name}, p)
	if repo.UpdatedAnnotationId != annotationId {
		t.Fatalf("expected to update annotation with id %v, got %v",
			annotationId, repo.UpdatedAnnotationId)
	}
	if repo.UpdatedLabelId != newLabel.Id {
		t.Fatalf("expected to update annotation with label id %v, got %v",
			newLabel.Id, repo.UpdatedLabelId)
	}
}
