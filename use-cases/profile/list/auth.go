package list

import (
	"context"
)

type Auth interface {
	ListUsers(ctx context.Context) error
}
