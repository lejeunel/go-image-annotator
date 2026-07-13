package create

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type Repo interface {
	Create(rl.Role) error
	Exists(string) (*bool, error)
}
