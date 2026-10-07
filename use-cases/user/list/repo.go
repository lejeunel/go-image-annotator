package list

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Store interface {
	List(pag.PaginationParams) ([]u.User, *pag.Pagination, error)
}
