package update

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type Repo interface {
	Update(rl.UpdatableModel) error
	Exists(string) (*bool, error)
}
