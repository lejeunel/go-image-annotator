package list

import (
	l "github.com/lejeunel/go-image-annotator/entities/label"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	List(pag.PaginationParams) ([]*l.Label, error)
	Count() (int64, error)
}
