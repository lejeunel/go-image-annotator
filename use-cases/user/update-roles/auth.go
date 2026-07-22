package update_role

import (
	"context"
)

type Auth interface {
	UpdateRoles(ctx context.Context) error
}
