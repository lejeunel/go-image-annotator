package update_group

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{},
		&fk.GroupRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := New(&fk.UserRepo{Missing: true},
		&fk.GroupRepo{ExistingNames: []string{"my-group"}})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Groups: []string{"my-group"}}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingGroupShouldFail(t *testing.T) {
	itr := New(&fk.UserRepo{}, &fk.GroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Groups: []string{"my-group"}}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := New(&fk.UserRepo{ErrOnFind: e.ErrInternal}, &fk.GroupRepo{ExistingNames: []string{"my-group"}})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAssignUserWhoIsAlreadyAssignedHasNoEffect(t *testing.T) {
	groups := []string{"a-group"}
	user := usr.NewUser("user@example.com",
		usr.WithGroups(groups))
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo, &fk.GroupRepo{ExistingNames: []string{"a-group"}})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Groups: []string{"a-group"}}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, groups, p.Got.Groups)
	assert.Nil(t, repo.GotNewGroup)
}

func TestUpdateGroups(t *testing.T) {
	user := usr.NewUser("user@example.com",
		usr.WithGroups([]string{"a-group"}))
	repo := &fk.UserRepo{Return: &user}
	itr := New(repo, &fk.GroupRepo{ExistingNames: []string{"a-group", "new-group", "another-new-group"}})
	p := &FakePresenter{}
	newGroups := []string{"new-group", "another-new-group"}
	itr.Execute(t.Context(), Request{Id: user.Id, Groups: newGroups}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, newGroups, repo.SetGroups_)
	assert.Equal(t, user.Id, repo.SetGroupsToUser)
}
