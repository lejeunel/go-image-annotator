package fetchall

import (
	"slices"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrWhenCountExceedsLimit(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Count_: 2}, WithLimit(1))
	itr.Execute(p)
	assert.ErrorIs(t, p.GotErr, e.ErrLabelLimitExceeded)
}

func TestHandleErrOnFetch(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFetch: true, Err: e.ErrInternal})
	itr.Execute(p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestFetchLabels(t *testing.T) {
	p := &FakePresenter{}
	labels := []string{"first-label", "second-labels"}
	itr := NewInteractor(&FakeRepo{Labels: labels})
	itr.Execute(p)
	assert.True(t, p.GotSuccess)
	assert.True(t, slices.Equal(p.Got.Labels, labels))
}
