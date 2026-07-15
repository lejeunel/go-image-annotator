package cloner

import (
	"fmt"
	task "github.com/lejeunel/go-image-annotator/entities/task"
)

type WorkerPool interface {
	Submit(func())
}

type Cloner struct {
	pool            WorkerPool
	taskFuncBuilder TaskFuncBuilder
}

func New(pool WorkerPool, taskFuncBuilder TaskFuncBuilder) Cloner {
	return Cloner{pool, taskFuncBuilder}
}

func (c Cloner) Clone(t task.CloneTask) error {
	errCtx := fmt.Errorf("cloning collection")
	fn, err := c.taskFuncBuilder.Build(t)
	if err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}
	c.pool.Submit(fn)
	return nil
}
