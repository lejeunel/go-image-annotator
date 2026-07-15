package cloner

import (
	task "github.com/lejeunel/go-image-annotator/entities/task"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskFuncBuilderFailure(t *testing.T) {
	taskFuncBuilder := &FakeTaskFuncBuilder{Err: e.ErrInternal}
	cloner := New(&FakeWorkerPool{}, taskFuncBuilder)
	err := cloner.Clone(task.NewCloneTask(task.NewTaskId(), "user@mail.com", "src", "dst"))
	assert.Error(t, err)
}

func TestWorkerPoolGetsCorrectTask(t *testing.T) {
	pool := &FakeWorkerPool{}
	cloner := New(pool, &FakeTaskFuncBuilder{Return: func() {}})
	task := task.NewCloneTask(task.NewTaskId(), "user@mail.com", "src", "dst")
	cloner.Clone(task)
	assert.NotNil(t, pool.Submitted)
}
