package update_group

import (
	"context"
)

type Auth interface {
	UpdateGroups(ctx context.Context) error
}
