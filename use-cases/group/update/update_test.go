package update

import (
	"testing"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.GroupRepo{}, WithAuth(fk.Auth{ErrOnAuth: e.ErrAuthorization}))
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
	group := g.NewGroup(g.NewGroupId(), "a-group")
	p := &FakePresenter{}
	repo := &fk.GroupRepo{Existing: []g.Group{group}}
	itr := New(repo)
	description := "updated-description"
	req := Request{
		Name:           group.Name,
		NewName:        "updated-name",
		NewDescription: &description,
	}
	itr.Execute(t.Context(), req, p)
	assert.Equal(t, req.NewName, p.Got.Name)
	assert.Equal(t, req.NewDescription, p.Got.Description)
}

func TestUpdateGroupWithNameAlreadyTakenShouldFail(t *testing.T) {
	p := &FakePresenter{}
	group := g.NewGroup(g.NewGroupId(), "a-group")
	anotherGroup := g.NewGroup(g.NewGroupId(), "another-group")
	itr := New(&fk.GroupRepo{Existing: []g.Group{group, anotherGroup}})
	itr.Execute(t.Context(), Request{Name: group.Name, NewName: anotherGroup.Name}, p)
	assert.True(t, p.GotDuplicationErr)
	assert.False(t, p.GotSuccess)
}

func TestUpdateGroupWithUnchangedNameShouldSucceed(t *testing.T) {
	p := &FakePresenter{}
	group := g.NewGroup(g.NewGroupId(), "a-group")
	itr := New(&fk.GroupRepo{Existing: []g.Group{group}})
	itr.Execute(t.Context(), Request{Name: group.Name, NewName: group.Name}, p)
	assert.True(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	group := g.NewGroup(g.NewGroupId(), "a-group")
	itr := New(&fk.GroupRepo{
		Existing:    []g.Group{group},
		ErrOnUpdate: e.ErrInternal,
	})
	itr.Execute(t.Context(),
		Request{Name: group.Name, NewName: group.Name}, p)
	assert.True(t, p.GotInternalErr)
}
