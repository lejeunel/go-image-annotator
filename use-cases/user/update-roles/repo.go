package update_role

import usr "github.com/lejeunel/go-image-annotator/entities/user"

type UserRepo interface {
	Find(string) (*usr.User, error)
	SetRoles(usr.UserId, []string) error
}

type RoleRepo interface {
	Exists(string) (*bool, error)
}
