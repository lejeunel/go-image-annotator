package find

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type Repo interface {
	Find(string) (*rl.Role, error)
}
