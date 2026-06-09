package list

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Repo interface {
	List(Request) ([]u.User, error)
	Count() (int64, error)
}
