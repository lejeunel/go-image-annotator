package list

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Response struct {
	Labels     []lbl.Label
	Pagination pagination.Pagination
}
