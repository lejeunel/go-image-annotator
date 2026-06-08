package delete

import (
	"context"
)

type Auth interface {
	DeleteUser(ctx context.Context) error
}
