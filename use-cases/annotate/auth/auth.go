package auth

import (
	"context"
)

type Auth interface {
	AnnotateGroup(ctx context.Context, group string) error
}
