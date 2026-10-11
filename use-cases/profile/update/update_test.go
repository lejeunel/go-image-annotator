package update

import (
	"testing"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile", pr.WithGroup("a-group"))
	itr := New(&fk.ProfileRepo{Existing: []pr.Profile{profile}},
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
	itr := New(&fk.ProfileRepo{Existing: []pr.Profile{}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: profile.Name}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestDestinationMustNotExist(t *testing.T) {
	p := &FakePresenter{}
	source := pr.NewProfile(pr.NewProfileId(), "source-profile")
	destination := pr.NewProfile(pr.NewProfileId(), "destination-profile")
	itr := New(
		&fk.ProfileRepo{Existing: []pr.Profile{source, destination}},
		&fk.GroupRepo{},
		&fk.LabelRepo{},
	)
	itr.Execute(t.Context(), Request{Name: "name", NewName: "new-name"}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestDestinationCanExistWhenUnchanged(t *testing.T) {
	p := &FakePresenter{}
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{Existing: []pr.Profile{profile}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: profile.Name, NewName: profile.Name}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleErrorOnLabelExist(t *testing.T) {
	p := &FakePresenter{}
	label := "a-label"
	profile := pr.NewProfile(pr.NewProfileId(), "a-profile")
	itr := New(&fk.ProfileRepo{Existing: []pr.Profile{profile}}, &fk.GroupRepo{},
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
	itr := New(&fk.ProfileRepo{Existing: []pr.Profile{profile}}, &fk.GroupRepo{},
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
	group := g.NewGroup(g.NewGroupId(), "a-group")
	newGroup := g.NewGroup(g.NewGroupId(), "new-group")
	profile := pr.NewProfile(pr.NewProfileId(), "profile-name", pr.WithGroup(group.Name))
	profileRepo := &fk.ProfileRepo{
		Existing: []pr.Profile{profile},
	}
	itr := New(profileRepo, &fk.GroupRepo{Existing: []g.Group{group, newGroup}}, &fk.LabelRepo{})
	itr.Execute(
		t.Context(),
		Request{Name: profile.Name, NewName: profile.Name, NewGroup: &newGroup.Name},
		p,
	)
	assert.NotNil(t, profileRepo.GotUpdateModel.NewGroup)
	assert.Equal(t, newGroup.Name, *profileRepo.GotUpdateModel.NewGroup)
	assert.True(t, p.GotSuccess)
}

func TestUpdateProfile(t *testing.T) {
	p := &FakePresenter{}
	newDescription := "new-description"
	group := g.NewGroup(g.NewGroupId(), "a-group")
	newGroup := g.NewGroup(g.NewGroupId(), "new-group")
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	current := pr.NewProfile(pr.NewProfileId(), "profile-name", pr.WithGroup(group.Name))
	req := Request{
		Name: current.Name, NewName: "new-profile-name",
		NewDescription: &newDescription, NewLabels: []string{label.Name},
		NewGroup: &newGroup.Name,
	}

	profileRepo := &fk.ProfileRepo{Existing: []pr.Profile{current}}

	itr := New(profileRepo, &fk.GroupRepo{Existing: []g.Group{newGroup}},
		&fk.LabelRepo{Existing: []lbl.Label{label}})
	itr.Execute(t.Context(), req, p)

	want := pr.UpdateModel{
		Name: current.Name, NewName: req.NewName, NewDescription: req.NewDescription,
		NewLabels: req.NewLabels, NewGroup: req.NewGroup,
	}
	assert.Equal(t, want, profileRepo.GotUpdateModel)
	assert.True(t, p.GotSuccess)
}
