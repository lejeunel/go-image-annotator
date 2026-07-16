package list

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Repo interface {
	List(pa.PaginationParams) ([]*clc.Collection, error)
	Count() (*int64, error)
}
