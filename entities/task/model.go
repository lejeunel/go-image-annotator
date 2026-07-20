package task

import (
	"fmt"
	e "github.com/lejeunel/go-image-annotator/entities/event"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type TaskType string

const (
	CollectionCloneTask TaskType = "collection-clone"
	IngestDirTask       TaskType = "ingest-dir"
)

func (r TaskType) String() string {
	return string(r)
}

func (r TaskType) Valid() bool {
	switch r {
	case CollectionCloneTask, IngestDirTask:
		return true
	default:
		return false
	}
}

func ParseTaskType(s string) (TaskType, error) {
	r := TaskType(s)
	if !r.Valid() {
		return "", fmt.Errorf("invalid task type %q", s)
	}
	return r, nil
}

type Task struct {
	Id     TaskId
	Type   TaskType
	Issuer u.UserId
	Events []e.Event
}

func NewTask(id TaskId, user u.UserId, type_ TaskType) Task {
	return Task{Id: id, Issuer: user,
		Type: type_}
}
