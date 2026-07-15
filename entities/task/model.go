package task

import (
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type TaskType int

const (
	CollectionCloneTask TaskType = iota
)

type TaskState int

const (
	PendingTask TaskState = iota
)

type TaskSpecs struct {
	Id     TaskId
	Issuer u.UserId
	Type   TaskType
	State  TaskState
}

func NewSpecs(id TaskId, user u.UserId, t TaskType) TaskSpecs {
	return TaskSpecs{Id: id, Issuer: user, Type: t, State: PendingTask}
}

type CloneTask struct {
	TaskSpecs
	Source      string
	Destination string
	Deep        bool
}

type Option func(*CloneTask)

func WithDeepClone() Option {
	return func(t *CloneTask) {
		t.Deep = true
	}
}

func NewCloneTask(id TaskId, user u.UserId, src, dst string, opts ...Option) CloneTask {
	cloneTask := &CloneTask{TaskSpecs: NewSpecs(id, user, CollectionCloneTask), Source: src, Destination: dst}

	for _, opt := range opts {
		opt(cloneTask)
	}
	return *cloneTask
}
