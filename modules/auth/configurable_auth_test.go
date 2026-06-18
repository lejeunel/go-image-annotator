package auth

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFailOnIllFormed(t *testing.T) {
	_, err := NewAuthRulesFromYaml(strings.NewReader("xy///"))
	assert.Error(t, err)
}

func TestFailOnNonExistingMethod(t *testing.T) {
	_, err := NewAuthRulesFromYaml(
		strings.NewReader(
			`
rules:
  - method: NonExistingMethod
`,
		))
	assert.Error(t, err)
}

func TestFailOnInvalidIgnoreGroupValue(t *testing.T) {
	_, err := NewAuthRulesFromYaml(
		strings.NewReader(
			`
rules:
  - method: CreateCollection
    ignore_group: maybe... I don't know
`,
		))
	assert.Error(t, err)
}

func TestFailOnInvalidRolesValue(t *testing.T) {
	_, err := NewAuthRulesFromYaml(
		strings.NewReader(
			`
rules:
  - method: CreateCollection
    ignore_group: maybe... I don't know
    roles: this-should-be-a-list
`,
		))
	assert.Error(t, err)
}

var validSpec = `
rules:
  - method: CreateCollection
    ignore_group: true
    roles: [a-role, another-role]
`

func TestValidRules(t *testing.T) {
	authRules, err := NewAuthRulesFromYaml(
		strings.NewReader(validSpec))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*authRules))
	assert.Equal(t, 2, len((*authRules)[0].Roles))
	assert.True(t, (*authRules)[0].IgnoreGroup)
}

func TestAuthConstruction(t *testing.T) {
	auth, err := NewConfigurableAuthFromYaml(strings.NewReader(validSpec))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(auth.Rules))
}

func TestAuthorizeWhenNoRuleSpecified(t *testing.T) {
	auth, err := NewConfigurableAuthFromYaml(strings.NewReader(""))
	assert.NoError(t, err)
	err = auth.CreateCollection(t.Context(), "")
	assert.NoError(t, err)
}

func TestNotAuthorizedWhenRequiredRoleIsMissing(t *testing.T) {
	auth, err := NewConfigurableAuth(
		[]AuthRule{{
			Method:      "CreateCollection",
			IgnoreGroup: true,
			Roles:       []string{"a-role-that-i-dont-have"}}})
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"my-role"}})
	err = auth.CreateCollection(ctx, "whatever")
	assert.Error(t, err)
}

func TestAuthorizedWhenRequiredRoleIsPresent(t *testing.T) {
	auth, err := NewConfigurableAuth(
		[]AuthRule{{
			Method:      "CreateCollection",
			IgnoreGroup: true,
			Roles:       []string{"a-role-that-i-have"}}})
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"a-role-that-i-have"}})
	err = auth.CreateCollection(ctx, "whatever")
	assert.NoError(t, err)
}

func TestNotAuthorizedWhenNotInGroup(t *testing.T) {
	auth, err := NewConfigurableAuth(
		[]AuthRule{{
			Method: "CreateCollection",
			Roles:  []string{"a-role-that-i-have"}}})
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"a-role-that-i-have"},
		Groups: []string{"group-of-losers"}})
	err = auth.CreateCollection(ctx, "group-of-chads")
	assert.Error(t, err)
}

func TestAuthorizedWhenMemberOfGroup(t *testing.T) {
	auth, err := NewConfigurableAuth(
		[]AuthRule{{
			Method: "CreateCollection",
			Roles:  []string{"a-role-that-i-have"}}})
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"a-role-that-i-have"},
		Groups: []string{"group-of-chads"}})
	err = auth.CreateCollection(ctx, "group-of-chads")
	assert.NoError(t, err)
}

func TestAuthorizedWhenNeededGroupIsVoid(t *testing.T) {
	auth, err := NewConfigurableAuth(
		[]AuthRule{{
			Method: "CreateCollection",
			Roles:  []string{"a-role-that-i-have"}}})
	assert.NoError(t, err)
	ctx := u.AppendUserToContext(t.Context(), u.User{Roles: []string{"a-role-that-i-have"}})
	err = auth.CreateCollection(ctx, "")
	assert.NoError(t, err)
}
