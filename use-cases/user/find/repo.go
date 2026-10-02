package find

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type UserStore interface {
	Find(u.UserId) (*u.User, error)
}
