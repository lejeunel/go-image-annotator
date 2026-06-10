package set_admin

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	SetAdmin(u.UserId, bool) error
}
