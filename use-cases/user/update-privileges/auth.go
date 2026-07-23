package update

import (
	"context"
)

type Auth interface {
	UpdateUserPrivileges(ctx context.Context) error
}
