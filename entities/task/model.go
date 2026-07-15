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

type CloneTask struct {
	TaskId
	Issuer      u.UserId
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
	cloneTask := &CloneTask{id, user, src, dst, false}

	for _, opt := range opts {
		opt(cloneTask)
	}
	return *cloneTask
}
