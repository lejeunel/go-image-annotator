package create

import (
	"testing"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(
		&fk.ProfileRepo{},
		&fk.LabelRepo{},
		&fk.GroupRepo{},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}),
	)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnCheckExistence(t *testing.T) {
	name := "my-profile"
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnExists: e.ErrInternal}, &fk.LabelRepo{},
		&fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateProfileWithDuplicateNameShouldFail(t *testing.T) {
	p := &FakePresenter{}
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}}, &fk.LabelRepo{},
		&fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: profile.Name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateWithInvalidNameShouldFail(t *testing.T) {
	name := "invalid-profile-name"
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{}, &fk.LabelRepo{}, &fk.GroupRepo{},
		WithNameValidator(&fk.StringValidator{Invalid: true}))
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestHandleInternalErrorOnCreate(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnCreate: e.ErrInternal},
		&fk.LabelRepo{}, &fk.GroupRepo{},
		WithNameValidator(&fk.StringValidator{}))
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestHandleErrorOnLabelExists(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{}, &fk.LabelRepo{ErrOnExists: e.ErrInternal},
		&fk.GroupRepo{})
	description := "a-description"
	req := Request{
		Name:        "a-profile",
		Description: &description,
		Labels:      []string{"the-label"},
	}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{}, &fk.LabelRepo{}, &fk.GroupRepo{})
	description := "a-description"
	req := Request{
		Name:        "a-profile",
		Description: &description,
		Labels:      []string{"the-label"},
	}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnGroupExists(t *testing.T) {
	p := &FakePresenter{}
	group := "the-group"
	itr := New(&fk.ProfileRepo{},
		&fk.LabelRepo{},
		&fk.GroupRepo{ErrOnExists: e.ErrInternal})
	req := Request{
		Name:  "a-profile",
		Group: &group,
	}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestCreate(t *testing.T) {
	p := &FakePresenter{}
	profileRepo := &fk.ProfileRepo{}
	group := "the-group"
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	labelRepo := &fk.LabelRepo{Existing: []lbl.Label{label}}
	itr := New(profileRepo, labelRepo, &fk.GroupRepo{ExistingNames: []string{group}})
	description := "a-description"
	req := Request{
		Name:        "a-profile",
		Description: &description,
		Labels:      []string{label.Name},
	}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, profileRepo.Created[0].Name, req.Name)
	assert.Equal(t, profileRepo.Created[0].Description, req.Description)
	assert.True(t, p.GotSuccess)
}
