package list

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
)

type Repo interface {
	List(Request) ([]*grp.Group, error)
	Count() (*int64, error)
}
