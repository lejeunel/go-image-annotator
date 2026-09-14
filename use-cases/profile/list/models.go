package list

import (
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Response struct {
	Profiles   []pr.Profile
	Pagination pagination.Pagination
}
