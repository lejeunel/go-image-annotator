package create

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(
		&fk.ProfileRepo{},
		&fk.LabelRepo{},
		&fk.GroupRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}),
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
	name := "my-profile"
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ExistingNames: []string{name}}, &fk.LabelRepo{},
		&fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: name}, p)
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
	req := Request{
		Name:        "a-profile",
		Description: "a-description",
		Labels:      []string{"the-label"},
	}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingLabelShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{}, &fk.LabelRepo{}, &fk.GroupRepo{})
	req := Request{
		Name:        "a-profile",
		Description: "a-description",
		Labels:      []string{"the-label"},
	}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnAddLabel(t *testing.T) {
	p := &FakePresenter{}
	labelName := "the-label"
	itr := New(&fk.ProfileRepo{ErrOnAddLabel: e.ErrInternal},
		&fk.LabelRepo{ExistingNames: []string{labelName}},
		&fk.GroupRepo{})
	req := Request{
		Name:        "a-profile",
		Description: "a-description",
		Labels:      []string{labelName},
	}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotInternalErr)
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
	labelName := "the-label"
	group := "the-group"
	labelRepo := &fk.LabelRepo{ExistingNames: []string{labelName}}
	itr := New(profileRepo, labelRepo, &fk.GroupRepo{ExistingNames: []string{group}})
	req := Request{
		Name:        "a-profile",
		Description: "a-description",
		Labels:      []string{labelName},
	}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, profileRepo.Created[0].Name, req.Name)
	assert.Equal(t, profileRepo.Created[0].Description, req.Description)
	assert.Equal(t, 1, len(profileRepo.AddedLabels))
	assert.Equal(t, labelName, profileRepo.AddedLabels[0])
	assert.True(t, p.GotSuccess)
}
