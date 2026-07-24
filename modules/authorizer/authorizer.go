package authorizer

//go:generate go run ./gen -struct Authorizer -in authorizer.go -out validmethods.gen.go -pkg authorizer

import (
	"context"
	"fmt"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"io"
)

type Authorizer struct {
	Rules map[string]AuthRule
}

func NewDefault() Authorizer {
	r := Authorizer{Rules: make(map[string]AuthRule)}
	return *r.SetAuthRules(DefaultRules)
}

func NewVoidAuth() Authorizer {
	return Authorizer{Rules: make(map[string]AuthRule)}
}

func New(rules []AuthRule) (*Authorizer, error) {
	ruleMap := make(map[string]AuthRule)
	for _, r := range rules {
		ruleMap[r.Method] = r
	}
	return &Authorizer{ruleMap}, nil
}

func NewFromYaml(r io.Reader) (*Authorizer, error) {
	rules, err := NewAuthRulesFromYaml(r)
	if err != nil {
		return nil, fmt.Errorf("building authorizer from yaml: %w", err)
	}
	return New(*rules)
}

func (a Authorizer) checkForGroup(userGroups []string, neededGroup string) error {
	for _, userGroup := range userGroups {
		if userGroup == neededGroup {
			return nil
		}
	}
	return fmt.Errorf("checking membership to group %v: %w", neededGroup, e.ErrAuthorization)
}
func (a Authorizer) checkForRole(userRoles, allowedRoles []string) error {

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
		return fmt.Errorf("checking for roles given allowed role set %v and assigned roles %v: %w", allowedRoles, userRoles, e.ErrAuthorization)
	}
	return nil
}
func (a Authorizer) check(ctx context.Context, method, group string) error {
	errCtx := "authorizing request"
	rule, ok := a.Rules[method]
	if !ok {
		return nil
	}
	user := u.IdentityFromContext(ctx)
	if user == nil {
		return fmt.Errorf("%v: fetching user info from context: %w", errCtx, e.ErrAuthorization)
	}

	if rule.AdminOnly {
		if user.IsAdmin() {
			return nil
		}
		return fmt.Errorf("%v: checking if user is admin given method is admin only: %w", errCtx, e.ErrAuthorization)
	}

	if user.IsAdmin() {
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
func (a *Authorizer) SetAuthRules(rules []AuthRule) *Authorizer {
	for _, r := range rules {
		a.Rules[r.Method] = r
	}
	return a
}
func (a Authorizer) CreateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "CreateCollection", group)
}
func (a Authorizer) DeleteCollection(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteCollection", group)
}
func (a Authorizer) UpdateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "UpdateCollection", group)
}
func (a Authorizer) CreateLabel(ctx context.Context) error {
	return a.check(ctx, "CreateLabel", "")
}
func (a Authorizer) DeleteLabel(ctx context.Context) error {
	return a.check(ctx, "DeleteLabel", "")
}
func (a Authorizer) UpdateLabel(ctx context.Context) error {
	return a.check(ctx, "UpdateLabel", "")
}
func (a Authorizer) Annotate(ctx context.Context, group string) error {
	return a.check(ctx, "Annotate", group)
}
func (a Authorizer) DeleteImage(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteImage", group)
}
func (a Authorizer) ImportImage(ctx context.Context, group string) error {
	return a.check(ctx, "ImportImage", group)
}
func (a Authorizer) IngestImage(ctx context.Context, group string) error {
	return a.check(ctx, "IngestImage", group)
}
func (a Authorizer) CreateUser(ctx context.Context) error {
	return a.check(ctx, "CreateUser", "")
}
func (a Authorizer) DeleteUser(ctx context.Context) error {
	return a.check(ctx, "DeleteUser", "")
}
func (a Authorizer) RenewToken(ctx context.Context) error {
	return a.check(ctx, "RenewToken", "")
}
func (a Authorizer) AssignUserToGroup(ctx context.Context) error {
	return a.check(ctx, "AssignUserToGroup", "")
}
func (a Authorizer) UnAssignUserFromGroup(ctx context.Context) error {
	return a.check(ctx, "UnAssignUserFromGroup", "")
}
func (a Authorizer) AssignRoleToUser(ctx context.Context) error {
	return a.check(ctx, "AssignRoleToUser", "")
}
func (a Authorizer) UnAssignRoleFromUser(ctx context.Context) error {
	return a.check(ctx, "UnAssignRoleFromUser", "")
}
func (a Authorizer) ListUsers(ctx context.Context) error {
	return a.check(ctx, "ListUsers", "")
}
func (a Authorizer) FindUser(ctx context.Context) error {
	return a.check(ctx, "FindUser", "")
}
func (a Authorizer) CreateGroup(ctx context.Context) error {
	return a.check(ctx, "CreateGroup", "")
}
func (a Authorizer) DeleteGroup(ctx context.Context) error {
	return a.check(ctx, "DeleteGroup", "")
}
func (a Authorizer) UpdateGroup(ctx context.Context) error {
	return a.check(ctx, "UpdateGroup", "")
}
func (a Authorizer) CreateRole(ctx context.Context) error {
	return a.check(ctx, "CreateRole", "")
}
func (a Authorizer) DeleteRole(ctx context.Context) error {
	return a.check(ctx, "DeleteRole", "")
}
func (a Authorizer) UpdateRole(ctx context.Context) error {
	return a.check(ctx, "UpdateRole", "")
}
func (a Authorizer) RequestForgottenPasswordToken(ctx context.Context) error {
	return a.check(ctx, "RequestForgottenPasswordToken", "")
}
func (a Authorizer) ChangePassword(ctx context.Context, id u.UserId) error {
	return a.check(ctx, "ChangePassword", "")
}
func (a Authorizer) CloneCollection(ctx context.Context, group string) error {
	return a.check(ctx, "CloneCollection", "")
}

func (a Authorizer) UpdateUserPrivileges(ctx context.Context) error {
	return a.check(ctx, "UpdateUserPrivileges", "")
}
