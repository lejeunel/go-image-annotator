package create

import (
	"context"
)

type Auth interface {
	CreateRole(ctx context.Context) error
}
