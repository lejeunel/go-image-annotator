package fetchall

import (
	"context"
)

type Auth interface {
	FetchAllLabels(ctx context.Context) error
}
