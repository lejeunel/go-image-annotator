package update

import (
	"testing"

	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.ProfileRepo{ReturnGroup: "a-group"}, &fk.GroupRepo{},
		&fk.LabelRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestHandleErrorOnCheckSourceExists(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnExists: e.ErrInternal}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{NewName: "new-name"}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestSourceMustExist(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "profile-name"}, p)
	assert.True(t, p.GotNotFoundErr)
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
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "name"}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleErrorOnLabelExist(t *testing.T) {
	p := &FakePresenter{}
	label := "a-label"
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{ErrOnExists: e.ErrInternal})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "name", NewLabels: []string{label}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestLabelMustExist(t *testing.T) {
	p := &FakePresenter{}
	label := "a-label"
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "name", NewLabels: []string{label}}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateGroup(t *testing.T) {
	p := &FakePresenter{}
	name := "profile-name"
	currentGroup := "current-group"
	newGroup := "new-group"
	profileRepo := &fk.ProfileRepo{ExistingNames: []string{name}, ReturnGroup: currentGroup}
	itr := New(profileRepo, &fk.GroupRepo{ExistingNames: []string{newGroup}}, &fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: name, NewName: name, NewGroup: &newGroup}, p)
	assert.NotNil(t, profileRepo.GotUpdateModel.NewGroup)
	assert.Equal(t, newGroup, *profileRepo.GotUpdateModel.NewGroup)
	assert.True(t, p.GotSuccess)
}

func TestUpdateProfile(t *testing.T) {
	p := &FakePresenter{}
	name := "profile-name"
	currentGroup := "current-group"

	newName := "new-profile-name"
	newGroup := "new-group"
	newDescription := "new-description"
	newLabels := []string{"new-label"}
	profileRepo := &fk.ProfileRepo{ExistingNames: []string{name}, ReturnGroup: currentGroup}
	itr := New(profileRepo, &fk.GroupRepo{ExistingNames: []string{newGroup}},
		&fk.LabelRepo{ExistingNames: newLabels})

	itr.Execute(t.Context(), Request{
		Name: name, NewName: newName,
		NewDescription: newDescription, NewLabels: newLabels,
		NewGroup: &newGroup,
	}, p)

	want := pr.UpdateModel{
		Name: name, NewName: newName, NewDescription: newDescription,
		NewLabels: newLabels, NewGroup: &newGroup,
	}
	assert.Equal(t, want, profileRepo.GotUpdateModel)
	assert.True(t, p.GotSuccess)
}
