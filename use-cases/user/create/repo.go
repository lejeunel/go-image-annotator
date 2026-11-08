package create

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	Create(usr.BaseUser) error
	Exists(string) (bool, error)
}

type GroupRepo interface {
	Exists(string) (*bool, error)
}

type RoleRepo interface {
	Exists(string) (*bool, error)
}
