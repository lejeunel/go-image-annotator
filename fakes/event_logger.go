package fake

import (
	"time"

	e "github.com/lejeunel/go-image-annotator/entities/event"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type EventLogger struct {
	ErrOnCount      error
	ErrOnList       error
	ErrOnFind       error
	ReturnTasks     []t.Task
	ReturnTask      t.Task
	InitializedTask bool
	Events          []e.Event
	Count_          int64
}

func (l *EventLogger) InitTask(t.TaskId, t.TaskType, u.UserId) error {
	l.InitializedTask = true
	return nil
}
func (l *EventLogger) FindTask(t.TaskId) (*t.Task, error) {
	if l.ErrOnFind != nil {
		return nil, l.ErrOnFind
	}
	return &l.ReturnTask, nil
}
func (l *EventLogger) AddEvent(t t.TaskId, e e.Event) error {
	l.Events = append(l.Events, e)
	return nil
}
func (l *EventLogger) Count(u.UserId) (*int64, error) {
	if l.ErrOnCount != nil {
		return nil, l.ErrOnCount
	}
	return &l.Count_, nil
}

func (l *EventLogger) ListUserTasks(user u.UserId, p pa.PaginationParams) ([]t.Task, error) {
	if l.ErrOnList != nil {
		return nil, l.ErrOnList
	}
	return l.ReturnTasks, nil
}

type EventLoggerRepo struct {
	ErrOnInitTask   error
	CreatedTaskId   t.TaskId
	CreatedTaskType t.TaskType
	CreatedTaskUser u.UserId
	CreatedTaskAt   time.Time

	ClippedTasksToNum int
	Count_            int64

	AddedEvents []e.Event
	ReturnTask  t.Task
}

func (l *EventLoggerRepo) FindTask(id t.TaskId) (*t.Task, error) {
	return &l.ReturnTask, nil
}

func (l *EventLoggerRepo) Count(user u.UserId) (*int64, error) {
	return &l.Count_, nil
}

func (l *EventLoggerRepo) CreateTask(
	taskId t.TaskId,
	now time.Time,
	taskType t.TaskType,
	user u.UserId,
) error {
	if l.ErrOnInitTask != nil {
		return l.ErrOnInitTask
	}
	l.CreatedTaskAt = now
	l.CreatedTaskId = taskId
	l.CreatedTaskType = taskType
	l.CreatedTaskUser = user
	return nil
}

func (l *EventLoggerRepo) ClipNumTasks(user u.UserId, n int) error {
	l.ClippedTasksToNum = n
	return nil
}

func (l *EventLoggerRepo) AddEvent(id t.TaskId, event e.Event) error {
	l.AddedEvents = append(l.AddedEvents, event)
	return nil
}

func (l *EventLoggerRepo) ListUserTasks(user u.UserId, p pa.PaginationParams) ([]t.Task, error) {
	return []t.Task{{
		Id:     l.CreatedTaskId,
		Type:   l.CreatedTaskType,
		Issuer: l.CreatedTaskUser,
	}}, nil
}
func (l *EventLoggerRepo) ClipTasks(u.UserId, int) error { return nil }
func (l *EventLoggerRepo) GetEvents(t.TaskId) ([]e.Event, error) {
	return l.AddedEvents, nil
}
