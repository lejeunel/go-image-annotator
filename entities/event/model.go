package event

import (
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	"time"
)

type Event struct {
	Timestamp time.Time
	TaskId    t.TaskId
	Issuer    u.UserId
	Type      t.TaskType
	State     t.TaskState
	Extra     map[string]string
	error     error
}

func (e *Event) SetError(err error) Event {
	e.error = err
	e.State = t.FailedTask
	e.Timestamp = time.Now()
	return *e
}

func (e *Event) SetStart() Event {
	e.State = t.StartedTask
	e.Timestamp = time.Now()
	return *e
}

func (e *Event) SetDone() Event {
	e.State = t.DoneTask
	e.Timestamp = time.Now()
	return *e
}

func New(taskId t.TaskId, issuer u.UserId, taskType t.TaskType) Event {
	return Event{
		Timestamp: time.Now(),
		TaskId:    taskId,
		Issuer:    issuer,
		Type:      taskType,
		State:     t.PendingTask}

}
