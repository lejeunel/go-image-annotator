package unassign_group

import (
	"context"
)

type Auth interface {
	UnAssignUserFromGroup(ctx context.Context, id string, group string) error
}
