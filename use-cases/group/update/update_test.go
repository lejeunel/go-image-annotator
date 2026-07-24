package update

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.GroupRepo{}, WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestUpdateNonExistingGroupShouldFail(t *testing.T) {
	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name, NewName: "new-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateGroup(t *testing.T) {
	name := "name"
	p := &FakePresenter{}
	repo := &fk.GroupRepo{ExistingNames: []string{name}}
	itr := New(repo)
	req := Request{Name: name,
		NewName:        "updated-name",
		NewDescription: "updated-description"}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.NewName, p.Got.Name)
	assert.Equal(t, req.NewDescription, p.Got.Description)
}

func TestUpdateGroupWithNameAlreadyTakenShouldFail(t *testing.T) {
	p := &FakePresenter{}
	name := "name"
	existing_name := "existing-name"
	itr := New(&fk.GroupRepo{ExistingNames: []string{name, existing_name}})
	itr.Execute(t.Context(), Request{Name: name, NewName: existing_name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateGroupWithUnchangedNameShouldSucceed(t *testing.T) {
	p := &FakePresenter{}
	name := "name"
	itr := New(&fk.GroupRepo{ExistingNames: []string{name}})
	itr.Execute(t.Context(), Request{Name: name, NewName: name}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	name := "name"
	itr := New(&fk.GroupRepo{ExistingNames: []string{name},
		ErrOnUpdate: e.ErrInternal})
	itr.Execute(t.Context(),
		Request{Name: name, NewName: name}, p)
	assert.True(t, p.GotInternalErr)
}
