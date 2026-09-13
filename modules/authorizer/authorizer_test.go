package authorizer

import (
	"strings"
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/stretchr/testify/assert"
)

func TestFailOnIllFormed(t *testing.T) {
	_, err := NewAuthRulesFromYaml(strings.NewReader("xy///"))
	assert.Error(t, err)
}

func TestFailOnNonExistingMethod(t *testing.T) {
	_, err := NewAuthRulesFromYaml(
		strings.NewReader(
			`
version: 1
rules:
  viewer:
    - NonExistingMethod
`,
		))
	assert.Error(t, err)
}

var validSpec = `
version: 1
rules:
  a-role:
    - CreateCollection
  another-role:
    - CreateCollection
  admin:
    - "*"
`

func TestValidRules(t *testing.T) {
	authRules, err := NewAuthRulesFromYaml(
		strings.NewReader(validSpec))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(*authRules))
}

func TestNotAuthorizedWhenRequiredRoleIsMissing(t *testing.T) {
	policies := map[string][]string{"super-role": {"CreatedCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"my-role"}})
	group := "whatever"
	err = auth.CreateCollection(ctx, &group)
	assert.Error(t, err)
}

func TestAuthorizedWhenRequiredRoleIsPresent(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(),
		u.User{Roles: []string{"a-role-that-i-have"}, Groups: []string{"my-group"}})
	group := "my-group"
	err = auth.CreateCollection(ctx, &group)
	assert.NoError(t, err)
}

func TestNotAuthorizedWhenNotInGroup(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{
		Roles:  []string{"a-role-that-i-have"},
		Groups: []string{"group-of-losers"},
	})
	group := "group-of-chads"
	err = auth.CreateCollection(ctx, &group)
	assert.Error(t, err)
}

func TestAuthorizedWhenMemberOfGroup(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{
		Roles:  []string{"a-role-that-i-have"},
		Groups: []string{"group-of-chads"},
	})
	group := "group-of-chads"
	err = auth.CreateCollection(ctx, &group)
	assert.NoError(t, err)
}

func TestAppendSetOfRules(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth := NewDefault()
	auth.SetAuthRules(policies)
	group := "the-group"
	err := auth.CreateCollection(t.Context(), &group)
	assert.Error(t, err)
}

func TestAdminDoesNotNeedRoleNorGroup(t *testing.T) {
	policies := map[string][]string{
		"a-role-that-i-dont-have": {"CreateCollection"},
		"admin":                   {"*"},
	}
	auth, _ := New(policies)
	ctx := u.AppendUserToContext(t.Context(),
		u.NewUser("admin@example.com", u.WithRoles([]string{"admin"})))
	group := "a-group-i-am-not-member-of"
	err := auth.Annotate(ctx, &group)
	assert.NoError(t, err)
}
