package read

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	FindUser(string) (*u.User, error)
}
