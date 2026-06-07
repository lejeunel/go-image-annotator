package create

import (
	"context"
)

type Auth interface {
	CreateUser(ctx context.Context) error
}
