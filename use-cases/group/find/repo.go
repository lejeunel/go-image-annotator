package find

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
)

type Repo interface {
	Find(string) (*grp.Group, error)
}
