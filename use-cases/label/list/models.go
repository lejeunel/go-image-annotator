package list

import (
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Request struct {
	Page     int64
	PageSize int
}

type LabelResponse struct {
	Name        string
	Description string
}

type Response struct {
	Labels     []LabelResponse
	Pagination pagination.Pagination
}
