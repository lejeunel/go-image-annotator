package unassign_group

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&FakeUserRepo{}, &FakeGroupRepo{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := New(&FakeUserRepo{Missing: true}, &FakeGroupRepo{})
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

func TestUnAssignUserWhoIsNotAssignedHasNoEffect(t *testing.T) {
	group := "a-group"
	user := usr.NewUser("user@example.com")
	repo := &FakeUserRepo{Return: &user}
	itr := New(repo, &FakeGroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: group}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Groups))
}

func TestUnAssignUser(t *testing.T) {
	group := "a-group"
	user := usr.NewUser("user@example.com",
		usr.WithGroups([]string{group}))
	repo := &FakeUserRepo{Return: &user}
	itr := New(repo, &FakeGroupRepo{})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: group}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Groups))
	assert.Equal(t, group, *repo.GotUnassignedGroup)
}
