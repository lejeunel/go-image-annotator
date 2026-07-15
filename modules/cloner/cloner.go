package cloner

import (
	"fmt"
	task "github.com/lejeunel/go-image-annotator/entities/task"
)

type WorkerPool interface {
	Submit(func())
}

type TaskFuncBuilder interface {
	Build(task.CloneTask) (func(), error)
}

type Cloner struct {
	WorkerPool
	TaskFuncBuilder
}

func New(pool WorkerPool, taskFuncBuilder TaskFuncBuilder) Cloner {
	return Cloner{pool, taskFuncBuilder}
}

func (c Cloner) Clone(t task.CloneTask) error {
	errCtx := fmt.Errorf("cloning collection")
	fn, err := c.TaskFuncBuilder.Build(t)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}
	c.WorkerPool.Submit(fn)
	return nil
}
