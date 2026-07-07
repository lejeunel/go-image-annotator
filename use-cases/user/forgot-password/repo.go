package forgot_password

import (
	"time"

	usr "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	AddForgottenPasswordState([]byte, usr.UserId, time.Time) error
	DeleteForgottenPasswordTokens(usr.UserId) error
	Exists(usr.UserId) (bool, error)
}
