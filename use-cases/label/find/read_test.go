package find

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnFind: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestReadLabel(t *testing.T) {
	label := l.NewLabel(l.NewLabelId(),
		"my-label",
		l.WithDescription("a-description"))
	repo := &fk.LabelRepo{Return: label}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), label.Name, p)
	assert.Equal(t, label, p.Got)
}
