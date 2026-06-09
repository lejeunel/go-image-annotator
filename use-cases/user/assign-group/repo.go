package assign_group

import usr "github.com/lejeunel/go-image-annotator/entities/user"

type UserRepo interface {
	Find(string) (*usr.User, error)
	AssignToGroup(string, string) error
}

type GroupRepo interface {
	Exists(string) error
}
