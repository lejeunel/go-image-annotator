package fake

import (
	"context"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Auth struct {
	Err error
}

func (f Auth) CreateCollection(ctx context.Context, g string) error {
	return f.Err
}

func (f Auth) UpdateCollection(ctx context.Context, g string) error {
	return f.Err
}

func (f Auth) CloneCollection(ctx context.Context, g string) error {
	return f.Err
}

func (f Auth) Annotate(ctx context.Context, g string) error {
	return f.Err
}

func (f Auth) DeleteImage(ctx context.Context, g string) error {
	return f.Err
}
func (f Auth) ImportImage(ctx context.Context, dstGroup string) error {
	return f.Err
}

func (f Auth) IngestImage(ctx context.Context, group string) error {
	return f.Err
}

func (f Auth) CreateLabel(ctx context.Context) error {
	return f.Err
}

func (f Auth) DeleteLabel(ctx context.Context) error {
	return f.Err
}

func (f Auth) FetchAllLabels(ctx context.Context) error {
	return f.Err
}

func (f Auth) ReadLabel(ctx context.Context) error {
	return f.Err
}

func (f Auth) ListLabels(ctx context.Context) error {
	return f.Err
}

func (f Auth) UpdateLabel(ctx context.Context) error {
	return f.Err
}

func (f Auth) CreateRole(ctx context.Context) error {
	return f.Err
}

func (f Auth) DeleteRole(ctx context.Context) error {
	return f.Err
}

func (f Auth) UpdateRole(ctx context.Context) error {
	return f.Err
}

func (f Auth) AssignUserToGroup(ctx context.Context) error {
	return f.Err
}

func (f Auth) AssignRoleToUser(ctx context.Context) error {
	return f.Err
}

func (f Auth) ChangePassword(ctx context.Context, id u.UserId) error {
	return f.Err
}

func (f Auth) CreateUser(ctx context.Context) error {
	return f.Err
}

func (f Auth) DeleteUser(ctx context.Context) error {
	return f.Err
}
