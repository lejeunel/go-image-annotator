package list

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	List(pag.PaginationParams) ([]u.User, error)
	Count() (int64, error)
}
