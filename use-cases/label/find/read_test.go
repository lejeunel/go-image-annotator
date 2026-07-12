package find

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadNonExistingLabelShouldFail(t *testing.T) {
	repo := &FakeRepo{Label: l.Label{Name: "my-label", Description: "a-description"}}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), "non-existing-label", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestReadLabel(t *testing.T) {
	label := l.NewLabel(l.NewLabelId(),
		"my-label",
		l.WithDescription("a-description"))
	repo := &FakeRepo{Label: label}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), label.Name, p)
	assert.Equal(t, label, p.Got)
}
