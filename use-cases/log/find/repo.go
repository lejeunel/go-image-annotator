package find

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
)

type TaskFinder interface {
	FindTask(t.TaskId) (*t.Task, error)
}
