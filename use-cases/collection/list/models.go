package list

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Response struct {
	Collections []clc.Collection
	Pagination  pagination.Pagination
}
