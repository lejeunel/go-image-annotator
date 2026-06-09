package update_label

import (
	"testing"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	lbl := lbl.NewLabel(lbl.NewLabelId(), "my-label")
	itr := New(&FakeRepo{Returns: &lbl},
		WithAuth(auth.FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(),
		Request{AnnotationId: a.NewAnnotationId().String()},
		p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrOnFindLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Err: e.ErrNotFound, ErrOnFindLabel: true})
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.ErrorIs(t, p.GotErr, e.ErrNotFound)
}

func TestHandleErrOnUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	itr := New(&FakeRepo{Returns: &label, Err: e.ErrNotFound, ErrOnUpdate: true})
	itr.Execute(t.Context(), Request{AnnotationId: a.NewAnnotationId().String()}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestFetchLabelFromName(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: &newLabel}
	itr := New(repo)
	itr.Execute(t.Context(), Request{AnnotationId: a.NewAnnotationId().String(), Label: newLabel.Name}, p)
	assert.Equal(t, repo.FetchedLabelWithName, newLabel.Name, "label name")
}

func TestUpdateLabel(t *testing.T) {
	p := &FakePresenter{}
	newLabel := lbl.NewLabel(lbl.NewLabelId(), "another-label")
	repo := &FakeRepo{Returns: &newLabel}
	itr := New(repo)
	annotationId := a.NewAnnotationId()
	itr.Execute(t.Context(), Request{AnnotationId: annotationId.String(), Label: newLabel.Name}, p)
	assert.Equal(t, repo.UpdatedAnnotationId, annotationId, "annotation id")
	assert.Equal(t, repo.UpdatedLabelId, newLabel.Id, "label id")
}
