package list

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Request struct {
	CollectionName *string
	Page           int64
	PageSize       int
}

type Response struct {
	Images     []im.Image
	Pagination pagination.Pagination
}
