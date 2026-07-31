package list

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Response struct {
	Tasks      []t.Task
	Pagination pagination.Pagination
}
