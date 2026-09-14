package find

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
)

type Repo interface {
	Find(pr.ProfileName) (*pr.Profile, error)
}
