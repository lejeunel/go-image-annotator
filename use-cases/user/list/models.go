package list

import (
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type UserResponse struct {
	Id     string
	Groups []string
	Roles  []string
}

type Response struct {
	Users      []UserResponse
	Pagination pagination.Pagination
}
