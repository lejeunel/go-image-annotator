package assign_group

import (
	"context"
)

type Auth interface {
	AssignUserToGroup(ctx context.Context, id string, group string) error
}
