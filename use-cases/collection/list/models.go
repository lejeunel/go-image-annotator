package list

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Request struct {
	PageSize int
	Page     int64
}

type Response struct {
	Collections []clc.Collection
	Pagination  pagination.Pagination
}
