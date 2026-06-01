package create

import (
	"context"
)

type Auth interface {
	CreateCollection(ctx context.Context, group string) error
}
