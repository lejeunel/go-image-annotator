package delete

import (
	"context"
)

type Auth interface {
	// delete collection given its group ownership
	DeleteGroup(ctx context.Context, group string) error
}
