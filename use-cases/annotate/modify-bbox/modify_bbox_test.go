package modify_bbox

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrNotFound, ErrOnFindLabel: true})
	itr.Execute(Request{}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected to get not found error")
	}
}

func TestInternalErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal, ErrOnFindLabel: true})
	itr.Execute(Request{}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestValidationErrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(Request{Xc: 1, Yc: 1, Width: -999, Height: 1}, p)
	if !p.GotValidationErr || p.GotSuccess {
		t.Fatalf("expected validation error")
	}
}

func TestNotFoundErrOnUpdateShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnUpdate: true, Err: e.ErrNotFound})
	itr.Execute(Request{Xc: 1, Yc: 1, Width: 1, Height: 1}, p)
	if !p.GotNotFoundErr || p.GotSuccess {
		t.Fatalf("expected not found error")
	}
}

func TestInternalErrOnUpdateShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnUpdate: true, Err: e.ErrInternal})
	itr.Execute(Request{Xc: 1, Yc: 1, Width: 1, Height: 1}, p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatalf("expected internal error")
	}
}

func TestUpdate(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	repo := &FakeRepo{Label: *label}
	itr := NewInteractor(repo)
	annotationId := a.NewAnnotationId()
	r := Request{AnnotationId: annotationId, Xc: 1, Yc: 1, Width: 1, Height: 1}
	itr.Execute(r, p)
	got := repo.Got
	want := a.BoundingBoxUpdatables{LabelId: label.Id, Xc: r.Xc,
		Yc: r.Yc, Width: r.Width, Height: r.Height}
	if !p.GotSuccess {
		t.Fatalf("expected success")
	}
	if got != want {
		t.Fatalf("expected to update with %+v, got %+v", want, got)
	}
}
