package delete

import (
	"context"
)

type Auth interface {
	DeleteCollection(ctx context.Context, group *string) error
}
