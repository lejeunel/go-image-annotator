package unassign_group

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.GroupRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := New(&fk.UserRepo{Missing: true}, &fk.GroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Group: "my-group"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingGroupShouldFail(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.GroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Group: "my-group"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := New(&fk.UserRepo{ErrOnFind: e.ErrInternal}, &fk.GroupRepo{ExistingNames: []string{"my-group"}})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Group: "my-group"}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestUnAssignUserWhoIsNotAssignedHasNoEffect(t *testing.T) {
	group := "a-group"
	user := usr.NewUser("user@example.com")
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo, &fk.GroupRepo{ExistingNames: []string{group}})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: group}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Groups))
}

func TestUnAssignUser(t *testing.T) {
	group := "a-group"
	user := usr.NewUser("user@example.com",
		usr.WithGroups([]string{group}))
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo, &fk.GroupRepo{ExistingNames: []string{group}})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: group}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Groups))
	assert.Equal(t, group, repo.GotUnassignedGroup)
}
