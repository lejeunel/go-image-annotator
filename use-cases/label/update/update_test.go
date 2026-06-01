package update

import (
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateNonExistingLabelShouldFail(t *testing.T) {

	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := NewInteractor(&FakeRepo{})
	itr.Execute(Request{Name: non_existing_name, NewName: "new-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateLabel(t *testing.T) {
	name := "name"

	p := &FakePresenter{}
	repo := &FakeRepo{Names: []string{name}}
	itr := NewInteractor(repo)
	req := Request{Name: name,
		NewName:        "updated-name",
		NewDescription: "updated-description"}
	itr.Execute(req, p)
	assert.Equal(t, p.Got.Name, req.NewName)
	assert.Equal(t, p.Got.Description, req.NewDescription)
}

func TestUpdateLabelWithNameAlreadyTakenShouldFail(t *testing.T) {

	p := &FakePresenter{}
	name := "name"
	existing_name := "existing-name"
	itr := NewInteractor(&FakeRepo{Names: []string{name, existing_name}})
	itr.Execute(Request{Name: name, NewName: existing_name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeErrRepo{e.ErrInternal})
	itr.Execute(Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}
