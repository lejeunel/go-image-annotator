package create

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.ProfileRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnCheckExistence(t *testing.T) {
	name := "my-profile"
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnExists: e.ErrInternal})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateProfileWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-profile"
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ExistingNames: []string{name}})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateWithInvalidNameShouldFail(t *testing.T) {
	name := "invalid-profile-name"
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{},
		WithNameValidator(&fk.StringValidator{Invalid: true}))
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestHandleInternalErrorOnCreate(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnCreate: e.ErrInternal},
		WithNameValidator(&fk.StringValidator{}))
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestCreate(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.ProfileRepo{}
	itr := New(repo)
	req := Request{
		Name:        "a-profile",
		Description: "a-description",
		Labels:      []string{"first-label", "second-label"}}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, repo.Created[0].Name, req.Name)
	assert.Equal(t, repo.Created[0].Description, req.Description)
	assert.Equal(t, repo.Created[0].Labels, req.Labels)
	assert.True(t, p.GotSuccess)
}
