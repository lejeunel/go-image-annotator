package delete

import (
	"testing"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(
		&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}),
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), profile.Name, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnIsUsed(t *testing.T) {
	p := &FakePresenter{}
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{
		ExistingProfiles: []pr.Profile{profile},
		ErrOnIsUsed:      e.ErrInternal,
	})
	itr.Execute(t.Context(), profile.Name, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteProfileWithAssociatedResourcesShouldFail(t *testing.T) {
	p := &FakePresenter{}
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{
		ExistingProfiles: []pr.Profile{profile},
		IsUsed_:          true,
	})

	itr.Execute(t.Context(), profile.Name, p)
	assert.True(t, p.GotDependencyErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalErrOnDelete(t *testing.T) {
	p := &FakePresenter{}
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{
		ExistingProfiles: []pr.Profile{profile},
		ErrOnDelete:      e.ErrInternal,
	})
	itr.Execute(t.Context(), profile.Name, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestDeleteProfile(t *testing.T) {
	p := &FakePresenter{}
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}})
	itr.Execute(t.Context(), profile.Name, p)
	assert.True(t, p.GotSuccess)
}
