package authorizer

import (
	"context"
)

type VoidAuthorizer struct{}

func NewVoidAuth() VoidAuthorizer {
	return VoidAuthorizer{}
}

func (a VoidAuthorizer) CreateCollection(ctx context.Context, group string) error {
	return nil
}
func (a VoidAuthorizer) DeleteCollection(ctx context.Context, group string) error {
	return nil
}
func (a VoidAuthorizer) UpdateCollection(ctx context.Context, group string) error {
	return nil
}
func (a VoidAuthorizer) CreateLabel(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) DeleteLabel(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) UpdateLabel(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) Annotate(ctx context.Context, group string) error {
	return nil
}
func (a VoidAuthorizer) DeleteImage(ctx context.Context, group string) error {
	return nil
}
func (a VoidAuthorizer) ImportImage(ctx context.Context, group string) error {
	return nil
}
func (a VoidAuthorizer) IngestImage(ctx context.Context, group string) error {
	return nil
}
func (a VoidAuthorizer) CreateUser(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) DeleteUser(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) AssignUserToGroup(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) UnAssignUserFromGroup(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) AssignRoleToUser(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) UnAssignRoleFromUser(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) ListUsers(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) FindUser(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) CreateGroup(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) DeleteGroup(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) UpdateGroup(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) CreateRole(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) DeleteRole(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) UpdateRole(ctx context.Context) error {
	return nil
}
func (a VoidAuthorizer) CloneCollection(ctx context.Context, group string) error {
	return nil
}

func (a VoidAuthorizer) UpdateUserPrivileges(ctx context.Context) error {
	return nil
}

func (a VoidAuthorizer) ReadPolicies(ctx context.Context) error {
	return nil
}
