package list

import (
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Request struct {
	PageSize int
	Page     int64
}

type Response struct {
	Groups     []*grp.Group
	Pagination pagination.Pagination
}
