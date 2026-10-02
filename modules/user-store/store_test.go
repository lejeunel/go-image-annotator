package userstore

import (
	"testing"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	r "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"

	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func TestFindUser(t *testing.T) {
	groups := []g.Group{g.NewGroup(g.NewGroupId(), "a-group")}
	roles := []r.Role{r.NewRole(r.NewRoleId(), "a-role")}
	user := u.NewUser("user@example.com", u.WithGroups(groups), u.WithRoles(roles))
	base := user.ToBase()
	store := UserStore{
		&fk.UserRepo{Return: &base},
		&fk.GroupRepo{Return: groups[0]},
		&fk.RoleRepo{Return: roles[0]},
	}

	r, err := store.Find(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, groups, r.Groups)
	assert.Equal(t, roles, r.Roles)
}
