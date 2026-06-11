package renew_token

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	SetAccessTokenHash(usr.UserId, []byte) error
	Exists(usr.UserId) (bool, error)
}
