package update

import (
	"context"
)

type Auth interface {
	UpdateCollection(context.Context, string) error
}
