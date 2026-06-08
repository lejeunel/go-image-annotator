package assign_group

import usr "github.com/lejeunel/go-image-annotator/entities/user"

type Repo interface {
	Find(string) (*usr.User, error)
	AssignToGroup(string, string) error
	UserExists(string) error
	GroupExists(string) error
}
