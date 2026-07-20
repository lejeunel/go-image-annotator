package event

import (
	"testing"
	"time"

	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	ur "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/user"
	ev "github.com/lejeunel/go-image-annotator/entities/event"
	ta "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestErrOnInitTaskShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteEventRepo(db)
	db.Close()
	err := repo.CreateTask(ta.NewTaskId(), time.Now(), ta.CollectionCloneTask, "user@example.com")
	assert.ErrorIs(t, err, e.ErrInternal)
}

func TestCreateAndRetrieveTaskOfUser(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteEventRepo(db)
	firstUser := u.NewUser("first")
	secondUser := u.NewUser("second")
	userRepo := ur.NewSQLiteUserRepo(db)
	userRepo.Create(firstUser)
	userRepo.Create(secondUser)
	repo.CreateTask(ta.NewTaskId(), time.Now(), ta.CollectionCloneTask, firstUser.Id)
	repo.CreateTask(ta.NewTaskId(), time.Now(), ta.CollectionCloneTask, secondUser.Id)
	tasks, err := repo.ListUserTasks(firstUser.Id, pa.PaginationParams{PageSize: 2, Page: 1})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, ta.CollectionCloneTask, tasks[0].Type)
}

func TestTasksAreRetrievedInInverseChronologicalOrder(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteEventRepo(db)
	user := u.NewUser("first")
	userRepo := ur.NewSQLiteUserRepo(db)
	userRepo.Create(user)
	t0 := ta.NewTaskId()
	t1 := ta.NewTaskId()
	repo.CreateTask(t0, time.Now(), ta.CollectionCloneTask, user.Id)
	repo.CreateTask(t1, time.Now(), ta.CollectionCloneTask, user.Id)
	tasks, _ := repo.ListUserTasks(user.Id, pa.PaginationParams{PageSize: 2, Page: 1})
	assert.Equal(t, tasks[0].Id, t1)
}

func TestAddingEventToNonExistingTaskShouldFail(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteEventRepo(db)
	err := repo.AddEvent(ta.NewTaskId(), ev.Event{})
	assert.Error(t, err)
}

func TestAddEventToTask(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteEventRepo(db)
	user := u.NewUser("first")
	userRepo := ur.NewSQLiteUserRepo(db)
	userRepo.Create(user)
	tid := ta.NewTaskId()
	repo.CreateTask(tid, time.Now(), ta.CollectionCloneTask, user.Id)
	event := ev.New()
	err := repo.AddEvent(tid, event.SetState(time.Now(), ev.StartedTask))
	assert.NoError(t, err)
	err = repo.AddEvent(tid, event.SetState(time.Now(), ev.FailedTask))
	assert.NoError(t, err)
	events, err := repo.GetEvents(tid)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, ev.FailedTask, events[0].State)
}

func TestClipNumTasksPerUser(t *testing.T) {
	db := s.NewInMemory()
	repo := NewSQLiteEventRepo(db)
	user := u.NewUser("first")
	userRepo := ur.NewSQLiteUserRepo(db)
	userRepo.Create(user)
	firstId := ta.NewTaskId()
	secondId := ta.NewTaskId()
	repo.CreateTask(firstId, time.Now(), ta.CollectionCloneTask, user.Id)
	repo.CreateTask(secondId, time.Now(), ta.CollectionCloneTask, user.Id)
	repo.AddEvent(firstId, ev.Event{})
	err := repo.ClipNumTasks(user.Id, 1)
	assert.NoError(t, err)
	r, _ := repo.ListUserTasks(user.Id, pa.PaginationParams{Page: 1, PageSize: 2})
	assert.Equal(t, 1, len(r))
	assert.Equal(t, secondId, r[0].Id)
}
