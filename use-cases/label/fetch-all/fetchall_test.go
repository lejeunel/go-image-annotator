package fetchall

import (
	"slices"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrWhenCountExceedsLimit(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{Count_: 2}, WithLimit(1))
	itr.Execute(t.Context(), p)
	assert.ErrorIs(t, p.GotErr, e.ErrLabelLimitExceeded)
}

func TestHandleErrOnFetch(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&FakeRepo{ErrOnFetch: true, Err: e.ErrInternal})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestFetchLabels(t *testing.T) {
	p := &FakePresenter{}
	labels := []string{"first-label", "second-labels"}
	itr := New(&FakeRepo{Labels: labels})
	itr.Execute(t.Context(), p)
	assert.True(t, p.GotSuccess)
	assert.True(t, slices.Equal(p.Got.Labels, labels))
}
