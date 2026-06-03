package read

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadNonExistingLabelShouldFail(t *testing.T) {
	repo := &FakeRepo{Label: l.Label{Name: "my-label", Description: "a-description"}}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	req := Request{Name: "non-existing-label"}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestReadLabel(t *testing.T) {
	label := l.NewLabel(l.NewLabelId(),
		"my-label",
		l.WithDescription("a-description"))
	repo := &FakeRepo{Label: label}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	req := Request{Name: label.Name}
	want := Response{Name: label.Name, Description: label.Description}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, want, p.Got)
}
