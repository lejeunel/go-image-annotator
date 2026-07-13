package update

import (
	"context"
)

type Auth interface {
	UpdateRole(ctx context.Context) error
}
