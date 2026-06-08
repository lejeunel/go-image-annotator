package auth

import (
	"context"
)

type PassThroughAuth struct {
}

func (a PassThroughAuth) CreateCollection(ctx context.Context, group string) error {
	return nil
}
func (a PassThroughAuth) DeleteCollection(ctx context.Context, group string) error {
	return nil
}
func (a PassThroughAuth) UpdateCollection(ctx context.Context, group string) error {
	return nil
}

func (a PassThroughAuth) CreateLabel(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) DeleteLabel(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) UpdateLabel(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) AnnotateGroup(ctx context.Context, group string) error {
	return nil
}

func (a PassThroughAuth) DeleteImage(ctx context.Context, group string) error {
	return nil
}

func (a PassThroughAuth) ImportImage(ctx context.Context, srcGroup string, dstGroup string) error {
	return nil
}

func (a PassThroughAuth) IngestImage(ctx context.Context, group string) error {
	return nil
}

func (a PassThroughAuth) CreateUser(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) DeleteUser(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) RenewToken(ctx context.Context, id string) error {
	return nil
}

func (a PassThroughAuth) AssignUserToGroup(ctx context.Context, id string, group string) error {
	return nil
}

func (a PassThroughAuth) UnAssignUserFromGroup(ctx context.Context, id string, group string) error {
	return nil
}

func (a PassThroughAuth) AssignRoleToUser(ctx context.Context, id string, role string) error {
	return nil
}
