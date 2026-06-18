package read

import (
	"context"
)

type Auth interface {
	FindUser(ctx context.Context) error
}
