package find

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	Find(u.UserId) (*u.User, error)
}
