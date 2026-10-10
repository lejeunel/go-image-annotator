package update

import (
	"testing"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.CollectionRepo{ReturnGroup: "a-group"}, &fk.GroupRepo{},
		&fk.ProfileRepo{},
		WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestUpdateNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	non_existing_name := "non-existing-name"
	itr := New(&fk.CollectionRepo{}, &fk.GroupRepo{}, &fk.ProfileRepo{})
	itr.Execute(t.Context(), Request{Name: non_existing_name, NewName: "new-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollection(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	p := &FakePresenter{}
	repo := &fk.CollectionRepo{
		Existing: []clc.Collection{collection},
	}
	itr := New(repo, &fk.GroupRepo{}, &fk.ProfileRepo{})
	description := "updated-description"
	req := Request{
		Name:           collection.Name,
		NewName:        "updated-name",
		NewDescription: &description,
	}
	itr.Execute(t.Context(), req, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, req.NewName, p.Got.Name)
	assert.Equal(t, req.NewDescription, p.Got.Description)
}

func TestUpdateCollectionWithNameAlreadyTakenShouldFail(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	anotherCollection := clc.NewCollection(clc.NewCollectionId(), "another-collection")
	itr := New(&fk.CollectionRepo{Existing: []clc.Collection{collection, anotherCollection}},
		&fk.GroupRepo{}, &fk.ProfileRepo{})
	itr.Execute(t.Context(), Request{Name: collection.Name, NewName: anotherCollection.Name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateCollectionWithNoGroup(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	itr := New(&fk.CollectionRepo{Existing: []clc.Collection{collection}},
		&fk.GroupRepo{}, &fk.ProfileRepo{})
	itr.Execute(t.Context(), Request{Name: collection.Name, NewName: collection.Name}, p)
	assert.True(t, p.GotSuccess)
}

func TestUpdateCollectionGroup(t *testing.T) {
	p := &FakePresenter{}
	currentGroup := "current-group"
	newGroup := "my-group"
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	clcRepo := &fk.CollectionRepo{Existing: []clc.Collection{collection}, ReturnGroup: currentGroup}
	itr := New(clcRepo, &fk.GroupRepo{ExistingNames: []string{newGroup}}, &fk.ProfileRepo{})
	itr.Execute(
		t.Context(),
		Request{Name: collection.Name, NewName: collection.Name, NewGroup: &newGroup},
		p,
	)
	assert.NotNil(t, clcRepo.GotUpdateModel.NewGroup)
	assert.Equal(t, newGroup, *clcRepo.GotUpdateModel.NewGroup)
	assert.True(t, p.GotSuccess)
}

func TestUpdateCollectionProfile(t *testing.T) {
	p := &FakePresenter{}
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	newName := "new-collection-name"
	currentProfile := "current-profile"
	newProfile := pr.NewProfile(pr.NewProfileId(), "new-profile")
	clcRepo := &fk.CollectionRepo{
		Existing:      []clc.Collection{collection},
		ReturnProfile: currentProfile,
	}
	itr := New(
		clcRepo,
		&fk.GroupRepo{},
		&fk.ProfileRepo{ExistingProfiles: []pr.Profile{newProfile}},
	)
	itr.Execute(
		t.Context(),
		Request{Name: collection.Name, NewName: newName, NewProfile: &newProfile.Name},
		p,
	)
	assert.NotNil(t, clcRepo.GotUpdateModel.NewProfile)
	assert.Equal(t, newProfile.Name, *clcRepo.GotUpdateModel.NewProfile)
	assert.True(t, p.GotSuccess)
}
