package unassign_group

import usr "github.com/lejeunel/go-image-annotator/entities/user"

type Repo interface {
	Find(string) (*usr.User, error)
	UnAssignFromGroup(string, string) error
	UserExists(string) error
	GroupExists(string) error
}
