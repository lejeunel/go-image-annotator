package read

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	Find(string) (*u.User, error)
}
