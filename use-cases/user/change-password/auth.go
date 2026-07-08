package change_password

import (
	"context"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Auth interface {
	ChangePassword(context.Context, u.UserId) error
}
