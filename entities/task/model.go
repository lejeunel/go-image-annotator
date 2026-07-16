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
	StartedTask
	FailedTask
	DoneTask
)

type CloneTask struct {
	TaskId
	Issuer      u.UserId
	Source      string
	Destination string
	Group       *string
	State       TaskState
	Deep        bool
}

func (t CloneTask) Type() TaskType {
	return CollectionCloneTask
}

type Option func(*CloneTask)

func WithDeepClone() Option {
	return func(t *CloneTask) {
		t.Deep = true
	}
}
func WithGroup(grp string) Option {
	return func(t *CloneTask) {
		t.Group = &grp
	}
}

func NewCloneTask(id TaskId, user u.UserId, src, dst string, opts ...Option) CloneTask {
	cloneTask := &CloneTask{TaskId: id, Issuer: user, Source: src, Destination: dst,
		State: PendingTask}

	for _, opt := range opts {
		opt(cloneTask)
	}
	return *cloneTask
}
