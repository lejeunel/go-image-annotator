package event_logger

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	ev "github.com/lejeunel/go-image-annotator/entities/event"
	ta "github.com/lejeunel/go-image-annotator/entities/task"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	repo := &fk.EventLoggerRepo{}
	logger := New(repo)
	now := time.Now()
	logger.Clock = clockwork.NewFakeClockAt(now)
	user := "user@example"
	taskId := ta.NewTaskId()
	taskType := ta.CollectionCloneTask
	logger.InitTask(taskId, taskType, user)
	assert.Equal(t, user, repo.CreatedTaskUser)
	assert.Equal(t, taskId, repo.CreatedTaskId)
	assert.Equal(t, taskType, repo.CreatedTaskType)
	assert.Equal(t, now, repo.CreatedTaskAt)
}

func TestCallClipTasksOnInit(t *testing.T) {
	repo := &fk.EventLoggerRepo{}
	logger := New(repo, WithMaxNumTasksPerUser(1))
	user := "user@example.com"
	logger.InitTask(ta.NewTaskId(), ta.CollectionCloneTask, user)
	assert.Equal(t, repo.ClippedTasksToNum, 1)
}

func TestAddEvent(t *testing.T) {
	repo := &fk.EventLoggerRepo{}
	logger := New(repo, WithMaxNumTasksPerUser(1))
	user := "user@example.com"
	id := ta.NewTaskId()
	now := time.Now()
	state := ev.FailedTask
	err := e.ErrDuplicate
	logger.InitTask(id, ta.CollectionCloneTask, user)
	logger.AddEvent(id, ev.Event{Time: now, State: state, Error: err.Error()})
	addedEvent := repo.AddedEvents[0]
	assert.Equal(t, now, addedEvent.Time)
	assert.Equal(t, state, addedEvent.State)
	assert.Equal(t, err.Error(), addedEvent.Error)
}
func TestRetrieveAggregate(t *testing.T) {
	repo := &fk.EventLoggerRepo{}
	logger := New(repo, WithMaxNumTasksPerUser(1))
	user := "user@example.com"
	id := ta.NewTaskId()
	now := time.Now()
	logger.InitTask(id, ta.CollectionCloneTask, user)
	logger.AddEvent(id, ev.Event{Time: now, State: ev.StartedTask})
	logger.AddEvent(id, ev.Event{Time: now.Add(time.Hour), State: ev.FailedTask, Error: e.ErrDuplicate.Error()})
	tasks, _ := logger.ListUserTasks(user, pa.PaginationParams{Page: 1, PageSize: 2})
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, 2, len(tasks[0].Events))
}
