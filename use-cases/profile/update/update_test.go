package update

import (
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.ProfileRepo{ReturnGroup: "a-group"}, &fk.GroupRepo{},
		&fk.LabelRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.False(t, p.GotSuccess)
	assert.True(t, p.GotAuthErr)
}

func TestHandleErrorOnCheckSourceExists(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ErrOnExists: e.ErrInternal}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{NewName: "new-name"}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestSourceMustExist(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "profile-name"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestDestinationMustNotExist(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name", "new-name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "new-name"}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

func TestDestinationCanExistWhenUnchanged(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "name"}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleErrorOnLabelExist(t *testing.T) {
	p := &FakePresenter{}
	label := "a-label"
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{ErrOnExists: e.ErrInternal})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "name", NewLabels: []string{label}}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestLabelMustExist(t *testing.T) {
	p := &FakePresenter{}
	label := "a-label"
	itr := New(&fk.ProfileRepo{ExistingNames: []string{"name"}}, &fk.GroupRepo{},
		&fk.LabelRepo{})
	itr.Execute(t.Context(), Request{Name: "name", NewName: "name", NewLabels: []string{label}}, p)
	assert.True(t, p.GotValidationErr)
	assert.False(t, p.GotSuccess)
}

// func TestUpdateCollection(t *testing.T) {
// 	name := "name"
// 	p := &FakePresenter{}
// 	repo := &fk.CollectionRepo{
// 		ExistingNames: []string{"name"},
// 		Return:        clc.NewCollection(clc.NewCollectionId(), name),
// 	}
// 	itr := New(repo, &fk.GroupRepo{})
// 	req := Request{
// 		Name:           name,
// 		NewName:        "updated-name",
// 		NewDescription: "updated-description",
// 	}
// 	itr.Execute(t.Context(), req, p)
// 	assert.True(t, p.GotSuccess)
// 	assert.Equal(t, req.NewName, p.Got.Name)
// 	assert.Equal(t, req.NewDescription, p.Got.Description)
// }

// func TestUpdateCollectionWithNameAlreadyTakenShouldFail(t *testing.T) {
// 	p := &FakePresenter{}
// 	name := "name"
// 	existing_name := "existing-name"
// 	itr := New(&fk.CollectionRepo{ExistingNames: []string{name, existing_name}},
// 		&fk.GroupRepo{})
// 	itr.Execute(t.Context(), Request{Name: name, NewName: existing_name}, p)
// 	assert.True(t, p.GotDuplicationErr)
// 	assert.False(t, p.GotSuccess)
// }

// func TestUpdateCollectionWithNoGroup(t *testing.T) {
// 	p := &FakePresenter{}
// 	name := "name"
// 	itr := New(&fk.CollectionRepo{ExistingNames: []string{name}, ErrOnGetGroup: e.ErrNotFound},
// 		&fk.GroupRepo{})
// 	itr.Execute(t.Context(), Request{Name: name, NewName: name}, p)
// 	assert.True(t, p.GotSuccess)
// }

// func TestUpdateCollectionGroup(t *testing.T) {
// 	p := &FakePresenter{}
// 	name := "name"
// 	currentGroup := "current-group"
// 	newGroup := "my-group"
// 	clcRepo := &fk.CollectionRepo{ExistingNames: []string{name}, ReturnGroup: currentGroup}
// 	itr := New(clcRepo, &fk.GroupRepo{ExistingNames: []string{newGroup}})
// 	itr.Execute(t.Context(), Request{Name: name, NewName: name, NewGroup: &newGroup}, p)
// 	assert.NotNil(t, clcRepo.GotUpdateModel.NewGroup)
// 	assert.Equal(t, newGroup, *clcRepo.GotUpdateModel.NewGroup)
// 	assert.True(t, p.GotSuccess)
// }
