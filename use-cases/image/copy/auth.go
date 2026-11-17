package copy

import (
	"context"
)

type Auth interface {
	CopyImage(ctx context.Context, group *string) error
}
