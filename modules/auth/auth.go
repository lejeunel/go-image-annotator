package auth

//go:generate go run ./gen -struct Auth -in auth.go -out validmethods.gen.go -pkg auth

import (
	"context"
	"fmt"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"io"
)

type Auth struct {
	Rules map[string]AuthRule
}

func NewDefault() Auth {
	r := Auth{Rules: make(map[string]AuthRule)}
	return *r.SetAuthRules(DefaultRules)
}

func NewVoidAuth() Auth {
	return Auth{Rules: make(map[string]AuthRule)}
}

func New(rules []AuthRule) (*Auth, error) {
	ruleMap := make(map[string]AuthRule)
	for _, r := range rules {
		ruleMap[r.Method] = r
	}
	return &Auth{ruleMap}, nil
}

func NewFromYaml(r io.Reader) (*Auth, error) {
	rules, err := NewAuthRulesFromYaml(r)
	if err != nil {
		return nil, fmt.Errorf("building authorizer from yaml: %w", err)
	}
	return New(*rules)
}

func (a Auth) checkForGroup(userGroups []string, neededGroup string) error {
	for _, userGroup := range userGroups {
		if userGroup == neededGroup {
			return nil
		}
	}
	return fmt.Errorf("checking membership to group %v: %w", neededGroup, e.ErrAuth)
}
func (a Auth) checkForRole(userRoles, allowedRoles []string) error {

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
		return fmt.Errorf("checking for roles given allowed role set %v and assigned roles %v: %w", allowedRoles, userRoles, e.ErrAuth)
	}
	return nil
}
func (a Auth) check(ctx context.Context, method, group string) error {
	errCtx := "authorizing request"
	rule, ok := a.Rules[method]
	if !ok {
		return nil
	}
	user := u.IdentityFromContext(ctx)
	if user == nil {
		return fmt.Errorf("%v: fetching user info from context: %w", errCtx, e.ErrAuth)
	}

	if rule.AdminOnly {
		if user.IsAdmin {
			return nil
		}
		return fmt.Errorf("%v: checking if user is admin given method is admin only: %w", errCtx, e.ErrAuth)
	}

	if user.IsAdmin {
		return nil
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
func (a *Auth) SetAuthRules(rules []AuthRule) *Auth {
	for _, r := range rules {
		a.Rules[r.Method] = r
	}
	return a
}

func (a Auth) CreateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "CreateCollection", group)
}
func (a Auth) DeleteCollection(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteCollection", group)
}
func (a Auth) UpdateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "UpdateCollection", group)
}
func (a Auth) CreateLabel(ctx context.Context) error {
	return a.check(ctx, "CreateLabel", "")
}
func (a Auth) DeleteLabel(ctx context.Context) error {
	return a.check(ctx, "DeleteLabel", "")
}
func (a Auth) UpdateLabel(ctx context.Context) error {
	return a.check(ctx, "UpdateLabel", "")
}
func (a Auth) Annotate(ctx context.Context, group string) error {
	return a.check(ctx, "Annotate", group)
}
func (a Auth) DeleteImage(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteImage", group)
}
func (a Auth) ImportImage(ctx context.Context, group string) error {
	return a.check(ctx, "ImportImage", group)
}
func (a Auth) IngestImage(ctx context.Context, group string) error {
	return a.check(ctx, "IngestImage", group)
}
func (a Auth) CreateUser(ctx context.Context) error {
	return a.check(ctx, "CreateUser", "")
}
func (a Auth) DeleteUser(ctx context.Context) error {
	return a.check(ctx, "DeleteUser", "")
}
func (a Auth) RenewToken(ctx context.Context) error {
	return a.check(ctx, "RenewToken", "")
}
func (a Auth) AssignUserToGroup(ctx context.Context) error {
	return a.check(ctx, "AssignUserToGroup", "")
}
func (a Auth) UnAssignUserFromGroup(ctx context.Context) error {
	return a.check(ctx, "UnAssignUserFromGroup", "")
}
func (a Auth) AssignRoleToUser(ctx context.Context) error {
	return a.check(ctx, "AssignRoleToUser", "")
}
func (a Auth) UnAssignRoleFromUser(ctx context.Context) error {
	return a.check(ctx, "UnAssignRoleFromUser", "")
}
func (a Auth) ListUsers(ctx context.Context) error {
	return a.check(ctx, "ListUsers", "")
}
func (a Auth) FindUser(ctx context.Context) error {
	return a.check(ctx, "FindUser", "")
}
func (a Auth) CreateGroup(ctx context.Context) error {
	return a.check(ctx, "CreateGroup", "")
}
func (a Auth) DeleteGroup(ctx context.Context) error {
	return a.check(ctx, "DeleteGroup", "")
}
func (a Auth) SetAdminRights(ctx context.Context) error {
	return a.check(ctx, "SetAdminRights", "")
}
