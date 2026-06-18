package auth

//go:generate go run ./gen -struct ConfigurableAuth -in configurable_auth.go -out validmethods.gen.go -pkg auth

import (
	"context"
	"fmt"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"io"
)

type ConfigurableAuth struct {
	Rules map[string]AuthRule
}

func NewVoidAuth() ConfigurableAuth {
	return ConfigurableAuth{}
}

func NewConfigurableAuth(rules []AuthRule) (*ConfigurableAuth, error) {
	ruleMap := make(map[string]AuthRule)
	for _, r := range rules {
		ruleMap[r.Method] = r
	}
	return &ConfigurableAuth{ruleMap}, nil
}

func NewConfigurableAuthFromYaml(r io.Reader) (*ConfigurableAuth, error) {
	rules, err := NewAuthRulesFromYaml(r)
	if err != nil {
		return nil, fmt.Errorf("building authorizer from yaml: %w", err)
	}

	return NewConfigurableAuth(*rules)

}
func (a ConfigurableAuth) checkForGroup(userGroups []string, neededGroup string) error {
	for _, userGroup := range userGroups {
		if userGroup == neededGroup {
			return nil
		}
	}
	return fmt.Errorf("checking membership to group %v: %w", neededGroup, e.ErrAuth)
}
func (a ConfigurableAuth) checkForRole(userRoles, allowedRoles []string) error {

	gotAdequateRole := false
	for _, userRole := range userRoles {
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				gotAdequateRole = true
				break
			}
		}
	}
	if !gotAdequateRole {
		return fmt.Errorf("checking for adequate role: %w", e.ErrAuth)
	}
	return nil
}
func (a ConfigurableAuth) check(ctx context.Context, method, group string) error {
	errCtx := "authorizing request"
	rule, ok := a.Rules[method]
	if !ok {
		return nil
	}
	user := u.IdentityFromContext(ctx)
	if user == nil {
		return fmt.Errorf("%v: fetching user info from context: %w", errCtx, e.ErrAuth)
	}
	if err := a.checkForRole(user.Roles, rule.Roles); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}
	if rule.IgnoreGroup || (group == "") {
		return nil
	}
	if err := a.checkForGroup(user.Groups, group); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}
	return nil
}

func (a ConfigurableAuth) CreateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "CreateCollection", group)
}
func (a ConfigurableAuth) DeleteCollection(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteCollection", group)
}
func (a ConfigurableAuth) UpdateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "UpdateCollection", group)
}
func (a ConfigurableAuth) CreateLabel(ctx context.Context) error {
	return a.check(ctx, "CreateLabel", "")
}
func (a ConfigurableAuth) DeleteLabel(ctx context.Context) error {
	return a.check(ctx, "DeleteLabel", "")
}
func (a ConfigurableAuth) UpdateLabel(ctx context.Context) error {
	return a.check(ctx, "UpdateLabel", "")
}
func (a ConfigurableAuth) AnnotateGroup(ctx context.Context, group string) error {
	return a.check(ctx, "AnnotateGroup", group)
}
func (a ConfigurableAuth) DeleteImage(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteImage", group)
}
func (a ConfigurableAuth) ImportImage(ctx context.Context, group string) error {
	return a.check(ctx, "ImportImage", group)
}
func (a ConfigurableAuth) IngestImage(ctx context.Context, group string) error {
	return a.check(ctx, "IngestImage", group)
}
func (a ConfigurableAuth) CreateUser(ctx context.Context) error {
	return a.check(ctx, "CreateUser", "")
}
func (a ConfigurableAuth) DeleteUser(ctx context.Context) error {
	return a.check(ctx, "DeleteUser", "")
}
func (a ConfigurableAuth) RenewToken(ctx context.Context) error {
	return a.check(ctx, "RenewToken", "")
}
func (a ConfigurableAuth) AssignUserToGroup(ctx context.Context) error {
	return a.check(ctx, "AssignUserToGroup", "")
}
func (a ConfigurableAuth) UnAssignUserFromGroup(ctx context.Context) error {
	return a.check(ctx, "UnAssignUserFromGroup", "")
}
func (a ConfigurableAuth) AssignRoleToUser(ctx context.Context) error {
	return a.check(ctx, "AssignRoleToUser", "")
}
func (a ConfigurableAuth) UnAssignRoleFromUser(ctx context.Context) error {
	return a.check(ctx, "UnAssignRoleFromUser", "")
}
func (a ConfigurableAuth) ListUsers(ctx context.Context) error {
	return a.check(ctx, "ListUsers", "")
}
func (a ConfigurableAuth) FindUser(ctx context.Context) error {
	return a.check(ctx, "FindUser", "")
}
func (a ConfigurableAuth) CreateGroup(ctx context.Context) error {
	return a.check(ctx, "CreateGroup", "")
}
func (a ConfigurableAuth) DeleteGroup(ctx context.Context) error {
	return a.check(ctx, "DeleteGroup", "")
}
func (a ConfigurableAuth) SetAdminRights(ctx context.Context) error {
	return a.check(ctx, "SetAdminRights", "")
}
