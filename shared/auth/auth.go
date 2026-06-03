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
