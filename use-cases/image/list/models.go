package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Request struct {
	im.FilterQueryStr
	pa.PaginationParams
	im.OrderingStr
}

type Response struct {
	Images     []im.Image
	Pagination pa.Pagination
}
