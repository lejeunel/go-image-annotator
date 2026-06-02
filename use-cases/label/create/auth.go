package create

import (
	"context"
)

type Auth interface {
	CreateLabel(ctx context.Context) error
}
