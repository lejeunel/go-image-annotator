package list

import (
	"context"
)

type Auth interface {
	// list collections
	ListCollection(ctx context.Context) error
}
