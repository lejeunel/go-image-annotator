package authorizer

//go:generate go run ./gen -struct Authorizer -in authorizer.go -out validmethods.gen.go -pkg authorizer

import (
	"context"
	"fmt"
	"io"
	"slices"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Authorizer struct {
	Rules Policies
}

func NewDefault() Authorizer {
	r := Authorizer{Rules: make(Policies)}
	return *r.SetAuthRules(DefaultPolicies)
}

func NewVoidAuth() Authorizer {
	return Authorizer{Rules: make(Policies)}
}

func New(rules Policies) (*Authorizer, error) {
	return &Authorizer{rules}, nil
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
func (a Authorizer) checkForRole(userRoles []string, method string) error {

	for _, role := range userRoles {
		if slices.Contains(a.Rules[role], method) {
			return nil
		}
		if slices.Contains(a.Rules[role], "*") {
			return nil
		}
	}
	return fmt.Errorf("checking for role access given user roles %v: %w", userRoles, e.ErrAuthorization)
}
func (a Authorizer) check(ctx context.Context, method string, group *string) error {
	errCtx := "authorizing request"
	user := u.IdentityFromContext(ctx)
	if user == nil {
		return fmt.Errorf("%v: fetching user info from context: %w", errCtx, e.ErrAuthorization)
	}
	if err := a.checkForRole(user.Roles, method); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}

	if (group != nil) && !user.IsAdmin() {
		if err := a.checkForGroup(user.Groups, *group); err != nil {
			return fmt.Errorf("%v: %w", errCtx, err)
		}

	}
	return nil
}
func (a *Authorizer) SetAuthRules(rules Policies) *Authorizer {
	a.Rules = rules
	return a
}
func (a Authorizer) CreateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "CreateCollection", &group)
}
func (a Authorizer) DeleteCollection(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteCollection", &group)
}
func (a Authorizer) UpdateCollection(ctx context.Context, group string) error {
	return a.check(ctx, "UpdateCollection", &group)
}
func (a Authorizer) CreateLabel(ctx context.Context) error {
	return a.check(ctx, "CreateLabel", nil)
}
func (a Authorizer) DeleteLabel(ctx context.Context) error {
	return a.check(ctx, "DeleteLabel", nil)
}
func (a Authorizer) UpdateLabel(ctx context.Context) error {
	return a.check(ctx, "UpdateLabel", nil)
}
func (a Authorizer) Annotate(ctx context.Context, group string) error {
	return a.check(ctx, "Annotate", &group)
}
func (a Authorizer) DeleteImage(ctx context.Context, group string) error {
	return a.check(ctx, "DeleteImage", &group)
}
func (a Authorizer) ImportImage(ctx context.Context, group string) error {
	return a.check(ctx, "ImportImage", &group)
}
func (a Authorizer) IngestImage(ctx context.Context, group string) error {
	return a.check(ctx, "IngestImage", &group)
}
func (a Authorizer) CreateUser(ctx context.Context) error {
	return a.check(ctx, "CreateUser", nil)
}
func (a Authorizer) DeleteUser(ctx context.Context) error {
	return a.check(ctx, "DeleteUser", nil)
}
func (a Authorizer) AssignUserToGroup(ctx context.Context) error {
	return a.check(ctx, "AssignUserToGroup", nil)
}
func (a Authorizer) UnAssignUserFromGroup(ctx context.Context) error {
	return a.check(ctx, "UnAssignUserFromGroup", nil)
}
func (a Authorizer) AssignRoleToUser(ctx context.Context) error {
	return a.check(ctx, "AssignRoleToUser", nil)
}
func (a Authorizer) UnAssignRoleFromUser(ctx context.Context) error {
	return a.check(ctx, "UnAssignRoleFromUser", nil)
}
func (a Authorizer) ListUsers(ctx context.Context) error {
	return a.check(ctx, "ListUsers", nil)
}
func (a Authorizer) FindUser(ctx context.Context) error {
	return a.check(ctx, "FindUser", nil)
}
func (a Authorizer) CreateGroup(ctx context.Context) error {
	return a.check(ctx, "CreateGroup", nil)
}
func (a Authorizer) DeleteGroup(ctx context.Context) error {
	return a.check(ctx, "DeleteGroup", nil)
}
func (a Authorizer) UpdateGroup(ctx context.Context) error {
	return a.check(ctx, "UpdateGroup", nil)
}
func (a Authorizer) CreateRole(ctx context.Context) error {
	return a.check(ctx, "CreateRole", nil)
}
func (a Authorizer) DeleteRole(ctx context.Context) error {
	return a.check(ctx, "DeleteRole", nil)
}
func (a Authorizer) UpdateRole(ctx context.Context) error {
	return a.check(ctx, "UpdateRole", nil)
}
func (a Authorizer) CloneCollection(ctx context.Context, group string) error {
	return a.check(ctx, "CloneCollection", nil)
}

func (a Authorizer) UpdateUserPrivileges(ctx context.Context) error {
	return a.check(ctx, "UpdateUserPrivileges", nil)
}
