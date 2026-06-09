package modify_bbox

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeRepo{},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String()},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Err: e.ErrNotFound, ErrOnFindLabel: true})
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String()},
		p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnFindLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Err: e.ErrInternal, ErrOnFindLabel: true})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestValidationErrShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{})
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String(), Xc: 1, Yc: 1, Width: -999, Height: 1}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestNotFoundErrOnUpdateShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnUpdate: true, Err: e.ErrNotFound})
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String(), Xc: 1, Yc: 1, Width: 1, Height: 1}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestInternalErrOnUpdateShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnUpdate: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{Xc: 1, Yc: 1, Width: 1, Height: 1}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdate(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	repo := &FakeRepo{Label: label}
	itr := New(repo)
	annotationId := a.NewAnnotationId()
	r := Request{AnnotationId: annotationId.String(), Xc: 1, Yc: 1, Width: 1, Height: 1}
	itr.Execute(t.Context(), r, p)
	got := repo.Got
	want := a.BoundingBoxUpdatables{LabelId: label.Id, Xc: r.Xc,
		Yc: r.Yc, Width: r.Width, Height: r.Height}
	assert.True(t, p.GotSuccess)
	if got != want {
		t.Fatalf("expected to update with %+v, got %+v", want, got)
	}
}
