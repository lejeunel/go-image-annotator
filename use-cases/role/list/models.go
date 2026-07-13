package list

import (
	rl "github.com/lejeunel/go-image-annotator/entities/role"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Request struct {
	PageSize int
	Page     int64
}

type Response struct {
	Roles      []*rl.Role
	Pagination pagination.Pagination
}
