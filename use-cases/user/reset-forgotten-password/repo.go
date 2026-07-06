package reset_forgotten_password

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	FindForgottenPassword([]byte) (*usr.ForgotPasswordState, error)
	UpdatePassword(usr.UserId, []byte) error
}
