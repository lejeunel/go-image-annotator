package renew_token

import (
	"context"
)

type Auth interface {
	RenewToken(ctx context.Context, id string) error
}
