package create

import (
	"context"
)

type Auth interface {
	CreateProfile(ctx context.Context) error
}
