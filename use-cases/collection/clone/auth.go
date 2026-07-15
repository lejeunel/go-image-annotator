package clone

import (
	"context"
)

type Auth interface {
	CloneCollection(ctx context.Context, group string) error
}
