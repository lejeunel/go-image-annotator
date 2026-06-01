package auth

import (
	"context"
	p "github.com/lejeunel/go-image-annotator/entities/principal"
)

type PrincipalProvider interface {
	Provide() (*p.Principal, error)
}

type PassThroughAuth struct {
}

func (a PassThroughAuth) CreateCollection(ctx context.Context, group string) error {
	return nil
}
func (a PassThroughAuth) DeleteCollection(ctx context.Context, group string) error {
	return nil
}
func (a PassThroughAuth) ListCollection(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) ReadCollection(context.Context) error {
	return nil
}

func (a PassThroughAuth) UpdateCollection(ctx context.Context, group string) error {
	return nil
}
