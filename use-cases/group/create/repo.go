package create

import (
	g "github.com/lejeunel/go-image-annotator/entities/group"
)

type Repo interface {
	Create(g.Group) error
	Exists(string) (*bool, error)
}
