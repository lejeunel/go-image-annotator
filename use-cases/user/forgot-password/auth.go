package forgot_password

import (
	"context"
)

type Auth interface {
	RequestForgottenPasswordToken(ctx context.Context) error
}
