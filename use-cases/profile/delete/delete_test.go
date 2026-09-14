package delete

import (
	"testing"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.ProfileRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnIsUsed(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnIsUsed: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteProfileWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{IsUsed_: true})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnDelete(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnDelete: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteProfile(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{Return: pr.NewProfile(pr.NewProfileId(), "my-profile")})
	itr.Execute(t.Context(), "my-profile", p)
	assert.True(t, p.GotSuccess)
}
