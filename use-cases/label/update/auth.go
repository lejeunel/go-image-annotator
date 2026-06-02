package update

import (
	"context"
)

type Auth interface {
	UpdateLabel(ctx context.Context) error
}
