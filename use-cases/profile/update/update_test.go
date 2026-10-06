package update

import (
	"testing"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile", pr.WithGroup("a-group"))
	itr := New(&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}},
		&fk.GroupRepo{},
		&fk.LabelRepo{},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Name: profile.Name}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestNonExistingSourceShouldFail(t *testing.T) {
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile", pr.WithGroup("a-group"))
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ExistingProfiles: []pr.Profile{}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: profile.Name}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestDestinationMustNotExist(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name", "new-name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "new-name"}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestDestinationCanExistWhenUnchanged(t *testing.T) {
	p := &FakePresenter{}
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: profile.Name, NewName: profile.Name}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleErrorOnLabelExist(t *testing.T) {
	p := &FakePresenter{}
	label := "a-label"
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}}, &fk.GroupRepo{},
		&fk.LabelRepo{ErrOnExists: e.ErrInternal})
	itr.Execute(
		t.Context(),
		Request{Name: profile.Name, NewName: profile.Name, NewLabels: []string{label}},
		p,
	)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestLabelMustExist(t *testing.T) {
	p := &FakePresenter{}
	label := "a-label"
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{ExistingProfiles: []pr.Profile{profile}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(
		t.Context(),
		Request{Name: profile.Name, NewName: profile.Name, NewLabels: []string{label}},
		p,
	)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateGroup(t *testing.T) {
	p := &FakePresenter{}
	currentGroup := "current-group"
	newGroup := "new-group"
	profile := pr.NewProfile(pr.NewProfileId(), "profile-name")
	profileRepo := &fk.ProfileRepo{
		ExistingProfiles: []pr.Profile{profile},
		ReturnGroup:      currentGroup,
	}
	itr := New(profileRepo, &fk.GroupRepo{ExistingNames: []string{newGroup}}, &fk.LabelRepo{})
	itr.Execute(
		t.Context(),
		Request{Name: profile.Name, NewName: profile.Name, NewGroup: &newGroup},
		p,
	)
	assert.NotNil(t, profileRepo.GotUpdateModel.NewGroup)
	assert.Equal(t, newGroup, *profileRepo.GotUpdateModel.NewGroup)
	assert.True(t, p.GotSuccess)
}

func TestUpdateProfile(t *testing.T) {
	p := &FakePresenter{}

	current := pr.NewProfile(pr.NewProfileId(), "profile-name", pr.WithGroup("current-group"))
	newDescription := "new-description"
	newGroup := "new-group"
	req := Request{
		Name: current.Name, NewName: "new-profile-name",
		NewDescription: &newDescription, NewLabels: []string{"new-label"},
		NewGroup: &newGroup,
	}

	profileRepo := &fk.ProfileRepo{ExistingProfiles: []pr.Profile{current}}
	itr := New(profileRepo, &fk.GroupRepo{ExistingNames: []string{newGroup}},
		&fk.LabelRepo{ExistingNames: req.NewLabels})
	itr.Execute(t.Context(), req, p)

	want := pr.UpdateModel{
		Name: current.Name, NewName: req.NewName, NewDescription: req.NewDescription,
		NewLabels: req.NewLabels, NewGroup: req.NewGroup,
	}
	assert.Equal(t, want, profileRepo.GotUpdateModel)
	assert.True(t, p.GotSuccess)
}
