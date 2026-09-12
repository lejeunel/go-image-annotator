package create

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type Repo interface {
	Create(pr.Profile) error
	Exists(string) (*bool, error)
}
