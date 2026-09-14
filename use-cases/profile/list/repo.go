package list

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	List(pag.PaginationParams) ([]pr.Profile, error)
	Count() (*int64, error)
}
