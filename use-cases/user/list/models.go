package list

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Response struct {
	Users      []u.User
	Pagination pagination.Pagination
}
