package delete

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type Repo interface {
	Delete(string) error
	IsUsed(string) (*bool, error)
	Find(string) (*pr.Profile, error)
}
