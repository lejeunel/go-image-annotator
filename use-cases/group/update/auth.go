package update

import (
	"context"
)

type Auth interface {
	UpdateGroup(context.Context) error
}
