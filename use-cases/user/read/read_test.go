package read

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"testing"

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

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Err: e.ErrInternal})
	itr.Execute(t.Context(), Request{}, p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestFindUser(t *testing.T) {
	groups := []string{"a-group"}
	roles := []string{"a-role"}
	user := u.NewUser("the-user-id",
		u.WithGroups(groups),
		u.WithRoles(roles))
	repo := &FakeRepo{Return: &user}
	p := &FakePresenter{}
	itr := NewInteractor(repo)
	req := Request{Id: user.Id}
	itr.Execute(t.Context(), req, p)

	want := Response{Id: user.Id, Groups: groups, Roles: roles}
	assert.Equal(t, want, p.Got)
}
