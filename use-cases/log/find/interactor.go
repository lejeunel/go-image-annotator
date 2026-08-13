package find

import (
	"context"
	"fmt"

	t "github.com/lejeunel/go-image-annotator/entities/task"
)

type Interactor struct {
	TaskFinder
}

func (i Interactor) Execute(ctx context.Context, id string, out OutputPort) {
	errCtx := "fetching task"

	taskId, err := t.NewTaskIdFromString(id)
	if err != nil {
		out.Error(fmt.Errorf("%v: parsing task id: %w", errCtx, err))
		return
	}

	found, err := i.TaskFinder.FindTask(*taskId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessFindTask(*found)
}

func New(r TaskFinder) Interactor {
	return Interactor{TaskFinder: r}
}
