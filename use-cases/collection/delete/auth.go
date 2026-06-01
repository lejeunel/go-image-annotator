package delete

import (
	"context"
)

type Auth interface {
	// delete collection given its group ownership
	DeleteCollection(ctx context.Context, group string) error
}
