package fetchall

import (
	"slices"
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnCount: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrWhenCountExceedsLimit(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{Count_: 2}, WithLimit(1))
	itr.Execute(t.Context(), p)
	assert.ErrorIs(t, p.GotErr, e.ErrLabelLimitExceeded)
}

func TestHandleErrOnFetch(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnFetch: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestFetchLabels(t *testing.T) {
	p := &FakePresenter{}
	labels := []string{"first-label", "second-labels"}
	itr := New(&fk.LabelRepo{ExistingNames: labels})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotSuccess)
	assert.True(t, slices.Equal(p.Got.Labels, labels))
}
