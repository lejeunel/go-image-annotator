package delete

import (
	"context"
)

type Auth interface {
	DeleteProfile(ctx context.Context, group *string) error
}
