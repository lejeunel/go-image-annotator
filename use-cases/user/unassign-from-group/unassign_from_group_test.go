package unassign_group

import (
	"testing"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := NewInteractor(&FakeRepo{},
		WithAuth(FailingAuth{}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingUserShouldFail(t *testing.T) {
	itr := NewInteractor(&FakeRepo{UserMissing: true})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Group: "my-group"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestMissingGroupShouldFail(t *testing.T) {
	itr := NewInteractor(&FakeRepo{GroupMissing: true})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: "user@example.com", Group: "my-group"}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleErrorOnFindUser(t *testing.T) {
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestUnAssignUserWhoIsNotAssignedHasNoEffect(t *testing.T) {
	group := "a-group"
	user := usr.NewUser("user@example.com")
	repo := &FakeRepo{Return: &user}
	itr := NewInteractor(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: group}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Groups))
	assert.Nil(t, repo.GotUnassignedGroup)
}

func TestUnAssignUser(t *testing.T) {
	group := "a-group"
	user := usr.NewUser("user@example.com",
		usr.WithGroups([]string{group}))
	repo := &FakeRepo{Return: &user}
	itr := NewInteractor(repo)
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{Id: user.Id, Group: group}, p)
	assert.True(t, p.GotSuccess)
	assert.Equal(t, 0, len(p.Got.Groups))
	assert.Equal(t, group, *repo.GotUnassignedGroup)
}
