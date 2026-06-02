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

func (a PassThroughAuth) CreateLabel(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) ListLabels(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) DeleteLabel(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) FetchAllLabels(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) ReadLabel(ctx context.Context) error {
	return nil
}

func (a PassThroughAuth) UpdateLabel(ctx context.Context) error {
	return nil
}
