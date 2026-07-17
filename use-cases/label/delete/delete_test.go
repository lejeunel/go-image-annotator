package delete

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.LabelRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteLabelWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{IsUsed_: true})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnIsUsed(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnIsUsed: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnExists(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnExists: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeletingMissingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{})
	itr.Execute(t.Context(), "a-label", p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteLabel(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ExistingNames: []string{"my-label"}})
	itr.Execute(t.Context(), "my-label", p)
	assert.True(t, p.GotSuccess)
}
