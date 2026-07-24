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
`

func TestValidRules(t *testing.T) {
	authRules, err := NewAuthRulesFromYaml(
		strings.NewReader(validSpec))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(*authRules))
}

func TestAuthConstruction(t *testing.T) {
	auth, err := NewFromYaml(strings.NewReader(validSpec))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(auth.Rules))
}

func TestNotAuthorizedWhenRequiredRoleIsMissing(t *testing.T) {
	policies := map[string][]string{"super-role": {"CreatedCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"my-role"}})
	err = auth.CreateCollection(ctx, "whatever")
	assert.Error(t, err)
}

func TestAuthorizedWhenRequiredRoleIsPresent(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(),
		u.User{Roles: []string{"a-role-that-i-have"}, Groups: []string{"my-group"}})
	err = auth.CreateCollection(ctx, "my-group")
	assert.NoError(t, err)
}

func TestNotAuthorizedWhenNotInGroup(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"a-role-that-i-have"},
		Groups: []string{"group-of-losers"}})
	err = auth.CreateCollection(ctx, "group-of-chads")
	assert.Error(t, err)
}

func TestAuthorizedWhenMemberOfGroup(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth, err := New(policies)
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"a-role-that-i-have"},
		Groups: []string{"group-of-chads"}})
	err = auth.CreateCollection(ctx, "group-of-chads")
	assert.NoError(t, err)
}

func TestAppendSetOfRules(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-have": {"CreateCollection"}}
	auth := NewVoidAuth()
	auth.SetAuthRules(policies)
	err := auth.CreateCollection(t.Context(), "")
	assert.Error(t, err)
}

func TestAdminDoesNotNeedRoleNorGroup(t *testing.T) {
	policies := map[string][]string{"a-role-that-i-dont-have": {"CreateCollection"},
		"admin": {"*"}}
	auth, _ := New(policies)
	ctx := u.AppendUserToContext(t.Context(),
		u.NewUser("admin@example.com", u.WithRoles([]string{"admin"})))
	err := auth.Annotate(ctx, "a-group-i-am-not-member-of")
	assert.NoError(t, err)
}
