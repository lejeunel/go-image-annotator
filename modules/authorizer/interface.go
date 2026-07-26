package authorizer

//go:generate go run ./gen -iface Interface -in interface.go -out validmethods.gen.go -pkg authorizer

import (
	"context"
)

type Interface interface {
	CreateCollection(ctx context.Context, group string) error
	DeleteCollection(ctx context.Context, group string) error
	UpdateCollection(ctx context.Context, group string) error
	CreateLabel(ctx context.Context) error
	DeleteLabel(ctx context.Context) error
	UpdateLabel(ctx context.Context) error
	Annotate(ctx context.Context, group string) error
	DeleteImage(ctx context.Context, group string) error
	ImportImage(ctx context.Context, group string) error
	IngestImage(ctx context.Context, group string) error
	CreateUser(ctx context.Context) error
	DeleteUser(ctx context.Context) error
	ListUsers(ctx context.Context) error
	FindUser(ctx context.Context) error
	CreateGroup(ctx context.Context) error
	DeleteGroup(ctx context.Context) error
	UpdateGroup(ctx context.Context) error
	CreateRole(ctx context.Context) error
	DeleteRole(ctx context.Context) error
	UpdateRole(ctx context.Context) error
	CloneCollection(ctx context.Context, group string) error
	UpdateUserPrivileges(ctx context.Context) error
}
