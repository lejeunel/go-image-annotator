package assign_role

import (
	"context"
)

type Auth interface {
	AssignRoleToUser(ctx context.Context, id string, role string) error
}
