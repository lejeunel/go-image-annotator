package authorizer

import (
	"testing"

	g "github.com/lejeunel/go-image-annotator/entities/group"
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	"github.com/stretchr/testify/assert"
)

func TestNotAuthorizedWhenRequiredMethodIsMissing(t *testing.T) {
	role := rl.NewRole(rl.NewRoleId(), "a-role-with-no-method")
	auth := New([]string{}, &fk.RoleRepo{Return: role})
	ctx := u.AppendUserToContext(t.Context(), u.User{})
	group := "whatever"
	err := auth.CreateCollection(ctx, &group)
	assert.Error(t, err)
}

func TestAuthorizedWhenRequiredRoleIsPresent(t *testing.T) {
	role := rl.NewRole(
		rl.NewRoleId(),
		"a-role-with-needed-method",
		rl.WithMethods([]string{"CreateCollection"}),
	)
	auth := New([]string{}, &fk.RoleRepo{Return: role})
	group := g.NewGroup(g.NewGroupId(), "my-group")
	ctx := u.AppendUserToContext(t.Context(),
		u.NewUser("user@example.com", u.WithRoles([]rl.Role{role}),
			u.WithGroups([]g.Group{group})))
	err := auth.CreateCollection(ctx, &group.Name)
	assert.NoError(t, err)
}

func TestNotAuthorizedWhenNotInGroup(t *testing.T) {
	role := rl.NewRole(rl.NewRoleId(), "a-role", rl.WithMethods([]string{"CreateCollection"}))
	user := u.NewUser("user@example.com", u.WithRoles([]rl.Role{role}))
	auth := New([]string{}, &fk.RoleRepo{Return: role})
	ctx := u.AppendUserToContext(t.Context(), user)
	group := "not-my-group"
	err := auth.CreateCollection(ctx, &group)
	assert.Error(t, err)
}

func TestAuthorizedWhenMemberOfGroup(t *testing.T) {
	role := rl.NewRole(rl.NewRoleId(), "a-role", rl.WithMethods([]string{"CreateCollection"}))
	group := g.NewGroup(g.NewGroupId(), "a-group")
	auth := New([]string{}, &fk.RoleRepo{Return: role})
	user := u.NewUser("user@example.com", u.WithRoles([]rl.Role{role}),
		u.WithGroups([]g.Group{group}))
	ctx := u.AppendUserToContext(t.Context(), user)

	err := auth.CreateCollection(ctx, &group.Name)
	assert.NoError(t, err)
}

func TestAdminDoesNotNeedRoleNorGroup(t *testing.T) {
	role := rl.NewRole(rl.NewRoleId(), "admin", rl.WithMethods([]string{"*"}))
	auth := New([]string{}, &fk.RoleRepo{Return: role})
	user := u.NewUser("admin@example.com", u.WithRoles([]rl.Role{role}))
	ctx := u.AppendUserToContext(t.Context(), user)
	group := "a-group-i-am-not-member-of"
	err := auth.Annotate(ctx, &group)
	assert.NoError(t, err)
}
