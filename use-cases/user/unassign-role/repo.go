package unassign_role

import usr "github.com/lejeunel/go-image-annotator/entities/user"

type Repo interface {
	Find(string) (*usr.User, error)
	UnAssignRole(string, string) error
	UserExists(string) error
}
