package list

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type Repo interface {
	List() ([]rl.Role, error)
}
