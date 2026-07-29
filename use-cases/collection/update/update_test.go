package update

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.CollectionRepo{ReturnGroup: "a-group"}, &fk.GroupRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestUpdateNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&fk.CollectionRepo{}, &fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name, NewName: "new-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollection(t *testing.T) {
	name := "name"
	p := &FakePresenter{}
	repo := &fk.CollectionRepo{ExistingNames: []string{"name"},
		Return: clc.NewCollection(clc.NewCollectionId(), name)}
	itr := New(repo, &fk.GroupRepo{})
	req := Request{Name: name,
		NewName:        "updated-name",
		NewDescription: "updated-description"}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, req.NewName, p.Got.Name)
	assert.Equal(t, req.NewDescription, p.Got.Description)
}

func TestUpdateCollectionWithNameAlreadyTakenShouldFail(t *testing.T) {

	p := &FakePresenter{}
	name := "name"
	existing_name := "existing-name"
	itr := New(&fk.CollectionRepo{ExistingNames: []string{name, existing_name}},
		&fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: name, NewName: existing_name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollectionWithNoGroup(t *testing.T) {
	p := &FakePresenter{}
	name := "name"
	itr := New(&fk.CollectionRepo{ExistingNames: []string{name}, ErrOnGetGroup: e.ErrNotFound},
		&fk.GroupRepo{})
	itr.Execute(t.Context(), Request{Name: name, NewName: name}, p)
	assert.True(t, p.GotSuccess)
}

func TestUpdateCollectionGroup(t *testing.T) {
	p := &FakePresenter{}
	name := "name"
	currentGroup := "current-group"
	newGroup := "my-group"
	clcRepo := &fk.CollectionRepo{ExistingNames: []string{name}, ReturnGroup: currentGroup}
	itr := New(clcRepo, &fk.GroupRepo{ExistingNames: []string{newGroup}})
	itr.Execute(t.Context(), Request{Name: name, NewName: name, NewGroup: &newGroup}, p)
	assert.NotNil(t, clcRepo.GotUpdateModel.NewGroup)
	assert.Equal(t, newGroup, *clcRepo.GotUpdateModel.NewGroup)
	assert.True(t, p.GotSuccess)
}
