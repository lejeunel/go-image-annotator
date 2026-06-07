package create

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	Create(usr.User) error
	Exists(string) (bool, error)
}
