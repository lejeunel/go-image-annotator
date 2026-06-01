package delete

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestDeleteLabelWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{IsUsed_: true})
	itr.Execute(Request{}, p)
	assert.Equal(t, true, p.GotDependencyErr)
	assert.Equal(t, false, p.GotSuccess)
}

func TestHandleInternalErrOnIsUsed(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal, ErrOnIsUsed: true})
	itr.Execute(Request{}, p)
	assert.Equal(t, true, p.GotInternalErr)
	assert.Equal(t, false, p.GotSuccess)
}

func TestHandleInternalErrOnExists(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal, ErrOnExists: true})
	itr.Execute(Request{}, p)
	assert.Equal(t, true, p.GotInternalErr)
	assert.Equal(t, false, p.GotSuccess)
}

func TestDeletingMissingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{IsMissing: true})
	itr.Execute(Request{}, p)
	assert.Equal(t, true, p.GotNotFoundErr)
	assert.Equal(t, false, p.GotSuccess)
}

func TestDeleteLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(Request{Name: "my-collection"}, p)
	assert.Equal(t, true, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(Request{}, p)
	assert.Equal(t, true, p.GotInternalErr)
	assert.Equal(t, false, p.GotSuccess)
}
