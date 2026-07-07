package reset_password

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	FindResetPasswordState([]byte) (*usr.ForgotPasswordState, error)
	UpdatePassword(usr.UserId, []byte) error
}
