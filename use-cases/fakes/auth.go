package fake

import (
	"context"
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
