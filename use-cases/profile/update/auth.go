package update

import (
	"context"
)

type Auth interface {
	UpdateProfile(context.Context, *string) error
}
