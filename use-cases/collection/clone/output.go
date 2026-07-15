package clone

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
)

type OutputPort interface {
	SuccessSubmitCloneTask(t.TaskSpecs)
	Error(error)
}
