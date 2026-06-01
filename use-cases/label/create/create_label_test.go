package create

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	v "github.com/lejeunel/go-image-annotator/shared/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateLabelWithDuplicateNameShouldFail(t *testing.T) {
	name := "my-label"
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Names: []string{name}})
	itr.Execute(Request{Name: name}, p)
	assert.Equal(t, true, p.GotDuplicationErr)
	assert.Equal(t, false, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(Request{Name: "a-name"}, p)
	assert.Equal(t, true, p.GotInternalErr)
}

func TestCreateLabelWithInvalidNameShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{}, WithNameValidator(&v.FakeNameValidator{Err: e.ErrValidation}))
	itr.Execute(Request{Name: "invalid-name"}, p)
	assert.Equal(t, true, p.GotValidationErr)
}

func TestCreateLabel(t *testing.T) {
	p := &FakePresenter{}
	repo := &FakeRepo{}
	itr := NewInteractor(repo)
	req := Request{Name: "a-name", Description: "a-description"}
	itr.Execute(req, p)

	assert.Equal(t, p.Got.Name, req.Name, "name")
	assert.Equal(t, p.Got.Description, req.Description, "description")
	assert.Equal(t, repo.Got.Name, req.Name, "name")
	assert.Equal(t, repo.Got.Description, req.Description, "description")
	assert.Equal(t, repo.Got.Id.IsNil(), false, "id")
}
