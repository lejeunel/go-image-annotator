package create

import (
	"context"
)

type Auth interface {
	CreateGroup(ctx context.Context) error
}
