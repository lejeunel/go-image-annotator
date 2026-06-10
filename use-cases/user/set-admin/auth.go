package set_admin

import (
	"context"
)

type Auth interface {
	SetAdminRights(ctx context.Context) error
}
