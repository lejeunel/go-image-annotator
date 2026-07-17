package list

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type LabelResponse struct {
	Name        string
	Description string
}

type Response struct {
	Labels     []lbl.Label
	Pagination pagination.Pagination
}
