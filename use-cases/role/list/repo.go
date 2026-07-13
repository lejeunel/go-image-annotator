package list

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
)

type Repo interface {
	List(Request) ([]*rl.Role, error)
	Count() (*int64, error)
}
