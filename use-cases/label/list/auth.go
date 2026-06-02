package list

import (
	"context"
)

type Auth interface {
	ListLabels(ctx context.Context) error
}
