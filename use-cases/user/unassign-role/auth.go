package unassign_role

import (
	"context"
)

type Auth interface {
	UnAssignRoleFromUser(ctx context.Context) error
}
