package fake

import (
	"time"

	e "github.com/lejeunel/go-image-annotator/entities/event"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type EventLogger struct {
}

func (l *EventLogger) InitTask(t.TaskId, t.TaskType, u.UserId, ...string) error { return nil }
func (l *EventLogger) AddEvent(time.Time, t.TaskId, e.State, error) error       { return nil }

type EventLoggerRepo struct {
	ErrOnInitTask   error
	CreatedTaskId   t.TaskId
	CreatedTaskType t.TaskType
	CreatedTaskUser u.UserId
	CreatedTaskAt   time.Time

	ClippedTasksToNum int

	AddedEvents []e.Event
}

func (l *EventLoggerRepo) CreateTask(taskId t.TaskId, now time.Time, taskType t.TaskType, user u.UserId) error {
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
