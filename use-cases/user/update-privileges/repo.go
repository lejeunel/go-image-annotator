package update

import usr "github.com/lejeunel/go-image-annotator/entities/user"

type UserRepo interface {
	Find(usr.UserId) (*usr.User, error)
	SetGroups(usr.UserId, []string) error
	SetRoles(usr.UserId, []string) error
}

type GroupRepo interface {
	Exists(string) (*bool, error)
}

type RoleRepo interface {
	Exists(string) (*bool, error)
}
