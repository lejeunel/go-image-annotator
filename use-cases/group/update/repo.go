package update

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
)

type Repo interface {
	Update(grp.UpdateModel) error
	Exists(string) (*bool, error)
}
