package create

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.RoleRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateRoleWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-role"
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ExistingNames: []string{name}})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ErrOnCreate: e.ErrInternal},
		WithNameValidator(&v.FakeNameValidator{}))
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
}

func TestCreateWithInvalidNameShouldFail(t *testing.T) {
	name := "my-role%/"
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ExistingNames: []string{name}},
		WithNameValidator(&v.FakeNameValidator{Err: e.ErrValidation}))
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreate(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.RoleRepo{}
	itr := New(repo)
	req := Request{Name: "a-role", Description: "a-description"}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, repo.Got.Name, req.Name)
	assert.Equal(t, repo.Got.Description, req.Description)
}
