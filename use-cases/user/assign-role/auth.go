package assign_role

import (
	"context"
)

type Auth interface {
	AssignRoleToUser(ctx context.Context) error
}
