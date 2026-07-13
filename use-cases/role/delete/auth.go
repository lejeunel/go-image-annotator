package delete

import (
	"context"
)

type Auth interface {
	DeleteRole(ctx context.Context) error
}
