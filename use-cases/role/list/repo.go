package list

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	List(pag.PaginationParams) ([]*rl.Role, error)
	Count() (*int64, error)
}
