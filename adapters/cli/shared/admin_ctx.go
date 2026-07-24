package shared

import (
	"context"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

func AnonymousAdminCtx() context.Context {
	return u.AppendUserToContext(context.Background(),
		u.NewUser("anonymous", u.WithRoles([]string{"admin"})))
}
