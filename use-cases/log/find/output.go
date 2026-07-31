package find

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
)

type OutputPort interface {
	SuccessFindTask(t.Task)
	Error(error)
}
