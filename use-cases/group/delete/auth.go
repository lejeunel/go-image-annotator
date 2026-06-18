package delete

import (
	"context"
)

type Auth interface {
	DeleteGroup(ctx context.Context) error
}
