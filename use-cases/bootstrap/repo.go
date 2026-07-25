package bootstrap

import (
	r "github.com/lejeunel/go-image-annotator/entities/role"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type UserRepo interface {
	Create(u.User) error
}

type RoleRepo interface {
	Create(r.Role) error
	Exists(string) (*bool, error)
}
