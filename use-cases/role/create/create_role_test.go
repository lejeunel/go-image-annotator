package create

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.RoleRepo{}, &fk.Auth{ErrOnAuth: e.ErrAuthorization})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateRoleWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-role"
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{ExistingNames: []string{name}}, &fk.Auth{})
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateWithInvalidNameShouldFail(t *testing.T) {
	name := "my-role%/"
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{}, &fk.Auth{},
		WithNameValidator(&fk.StringValidator{Invalid: true}))
	itr.Execute(t.Context(), Request{Name: name}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreateWithInvalidMethodShouldFail(t *testing.T) {
	name := "my-role"
	p := &FakePresenter{}
	itr := New(&fk.RoleRepo{}, &fk.Auth{ExistingMethods: []string{"a-method"}})
	itr.Execute(t.Context(), Request{Name: name, Methods: []string{"a-non-existing-method"}}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreate(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.RoleRepo{}
	itr := New(repo, &fk.Auth{ExistingMethods: []string{"a-method"}})
	methods := []string{"a-method"}
	description := "a-description"
	req := Request{Name: "a-role", Description: &description, Methods: methods}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, repo.Created[0].Name, req.Name)
	assert.Equal(t, repo.Created[0].Description, req.Description)
	assert.Equal(t, repo.Created[0].Methods, req.Methods)
}
