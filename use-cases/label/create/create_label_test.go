package create

import (
	"testing"

	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.LabelRepo{}, WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestCreateLabelWithDuplicateNameShouldFail(t *testing.T) {
	p := &FakePresenter{}
	label := lbl.NewLabel(lbl.NewLabelId(), "a-label")
	itr := New(&fk.LabelRepo{Existing: []lbl.Label{label}})
	itr.Execute(t.Context(), Request{Name: label.Name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{ErrOnCreate: e.ErrInternal})
	itr.Execute(t.Context(), Request{Name: "a-name"}, p)
	assert.True(t, p.GotInternalErr)
}

func TestCreateLabelWithInvalidNameShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.LabelRepo{}, WithNameValidator(&fk.StringValidator{Invalid: true}))
	itr.Execute(t.Context(), Request{Name: "invalid-name"}, p)
	assert.True(t, p.GotValidationErr)
}

func TestCreateLabel(t *testing.T) {
	p := &FakePresenter{}
	repo := &fk.LabelRepo{}
	itr := New(repo)
	description := "a-description"
	req := Request{Name: "a-name", Description: &description}
	itr.Execute(t.Context(), req, p)

	assert.Equal(t, p.Got.Name, req.Name)
	assert.Equal(t, p.Got.Description, req.Description)
	assert.Equal(t, repo.Created.Name, req.Name)
	assert.Equal(t, repo.Created.Description, req.Description)
	assert.False(t, repo.Created.Id.IsNil())
}
