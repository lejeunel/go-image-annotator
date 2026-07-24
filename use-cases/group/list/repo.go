package list

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
)

type Repo interface {
	List() ([]grp.Group, error)
}
