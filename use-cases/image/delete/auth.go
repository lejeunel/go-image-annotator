package delete

import (
	"context"
)

type Auth interface {
	DeleteImage(ctx context.Context, group string) error
}
