package assign_role

import usr "github.com/lejeunel/go-image-annotator/entities/user"

type UserRepo interface {
	Find(string) (*usr.User, error)
	AssignRole(string, string) error
}
