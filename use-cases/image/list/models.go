package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Request struct {
	im.FilteringParams
	im.OrderingParams
}

type Response struct {
	Images     []im.Image
	Pagination pagination.Pagination
}
