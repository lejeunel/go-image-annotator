package event_logger

import (
	"fmt"

	"github.com/jonboulle/clockwork"
	e "github.com/lejeunel/go-image-annotator/entities/event"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"time"
)

type Interface interface {
	InitTask(t.TaskId, t.TaskType, u.UserId, ...string) error
	AddEvent(time.Time, t.TaskId, e.State, error) error
}

type Repo interface {
	CreateTask(t.TaskId, time.Time, t.TaskType, u.UserId) error
	ListUserTasks(u.UserId, pa.PaginationParams) ([]t.Task, error)
	AddEvent(t.TaskId, e.Event) error
	GetEvents(t.TaskId) ([]e.Event, error)
	ClipNumTasks(u.UserId, int) error
}

type EventLogger struct {
	Repo
	clockwork.Clock
	clipNumTasks *int
}

func (l EventLogger) InitTask(id t.TaskId, type_ t.TaskType, user u.UserId) error {
	errCtx := fmt.Errorf("initializing task with id %v, type %v for user %v", id, type_, user)
	if err := l.Repo.CreateTask(id, l.Clock.Now(), type_, user); err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}

	if l.clipNumTasks != nil {
		err := l.ClipNumTasks(user, *l.clipNumTasks)
		if err != nil {
			return fmt.Errorf("%w: clipping oldest tasks: %w", errCtx, err)
		}
	}
	return nil
}
func (l EventLogger) ListUserTasks(user u.UserId, p pa.PaginationParams) ([]t.Task, error) {
	errCtx := fmt.Errorf("listing user tasks for %v", user)
	tasks, err := l.Repo.ListUserTasks(user, p)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}

	for i, _ := range tasks {
		events, err := l.Repo.GetEvents(tasks[i].Id)
		if err != nil {
			return nil, fmt.Errorf("%w: retrieving events: %w", errCtx, err)
		}
		tasks[i].Events = events
	}
	return tasks, nil
}

type Option func(*EventLogger)

func WithMaxNumTasksPerUser(nTasks int) Option {
	return func(i *EventLogger) {
		i.clipNumTasks = &nTasks
	}
}

func New(r Repo, opts ...Option) EventLogger {
	l := &EventLogger{Repo: r,
		Clock: clockwork.NewRealClock(),
	}
	for _, opt := range opts {
		opt(l)
	}
	return *l
}
