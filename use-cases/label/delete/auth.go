package delete

import (
	"context"
)

type Auth interface {
	DeleteLabel(ctx context.Context) error
}
