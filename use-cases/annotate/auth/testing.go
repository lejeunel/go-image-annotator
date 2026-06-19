package auth

import (
	"context"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type FailingAuth struct {
}

func (f FailingAuth) Annotate(ctx context.Context, g string) error {
	return e.ErrAuth
}
