package assign_group

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeUserRepo{},
		&FakeGroupRepo{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := New(&FakeUserRepo{Missing: true},
		&FakeGroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Group: "my-group"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingGroupShouldFail(t *testing.T) {
	itr := New(&FakeUserRepo{}, &FakeGroupRepo{Missing: true})

	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Group: "my-group"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := New(&FakeUserRepo{Err: e.ErrInternal}, &FakeGroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestAssignUserWhoIsAlreadyAssignedHasNoEffect(t *testing.T) {
	groups := []string{"a-group"}
	user := usr.NewUser("user@example.com",
		usr.WithGroups(groups))
	repo := &FakeUserRepo{Return: &user}
	itr := New(repo, &FakeGroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: "a-group"}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, groups, p.Got.Groups)
	assert.Nil(t, repo.GotNewGroup)
}

func TestAssignUser(t *testing.T) {
	user := usr.NewUser("user@example.com",
		usr.WithGroups([]string{"a-group"}))
	newGroup := "new-group"
	updatedGroups := []string{"a-group", newGroup}
	repo := &FakeUserRepo{Return: &user}
	itr := New(repo, &FakeGroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: newGroup}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, updatedGroups, p.Got.Groups)
	assert.Equal(t, newGroup, *repo.GotNewGroup)
}
