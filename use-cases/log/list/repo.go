package list

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type TaskLogger interface {
	ListUserTasks(user u.UserId, p pa.PaginationParams) ([]t.Task, error)
	Count(user u.UserId) (*int64, error)
}
