package auth

import (
	"context"
)

type Auth interface {
	Annotate(ctx context.Context, group string) error
}
