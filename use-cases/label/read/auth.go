package read

import (
	"context"
)

type Auth interface {
	ReadLabel(ctx context.Context) error
}
