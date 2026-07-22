package find

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"testing"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	itr := New(&fk.UserRepo{},
		WithAuth(fk.Auth{Err: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestHandleInternalError(t *testing.T) {
	p := &FakePresenter{}
	itr := New(&fk.UserRepo{ErrOnFind: e.ErrInternal})
	itr.Execute(t.Context(), "", p)
	assert.True(t, p.GotInternalErr)
	assert.False(t, p.GotSuccess)
}

func TestFindUser(t *testing.T) {
	groups := []string{"a-group"}
	roles := []string{"a-role"}
	user := u.NewUser("the-user-id",
		u.WithGroups(groups),
		u.WithRoles(roles))
	repo := &fk.UserRepo{Return: &user}
	p := &FakePresenter{}
	itr := New(repo)
	itr.Execute(t.Context(), user.Id, p)

	want := Response{Id: user.Id, Groups: groups, Roles: roles}
	assert.Equal(t, want, p.Got)
}
