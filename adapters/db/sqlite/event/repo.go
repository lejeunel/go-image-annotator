package event

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	ev "github.com/lejeunel/go-image-annotator/entities/event"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type SQLiteEventRepo struct {
	Db adb.Querier
}

// ClipNumTasksPerUser(int) error

type Task struct {
	Id        t.TaskId   `db:"id"`
	User      u.UserId   `db:"user_id"`
	CreatedAt time.Time  `db:"created_at"`
	Type      t.TaskType `db:"type_"`
}

type Event struct {
	Id    t.TaskId  `db:"task_id"`
	Time  time.Time `db:"time"`
	Error string    `db:"error"`
	State ev.State  `db:"state"`
	Extra string    `db:"extra"`
}

func (r SQLiteEventRepo) CreateTask(id t.TaskId, now time.Time, taskType t.TaskType, user u.UserId) error {
	query := `INSERT INTO tasks (id, user_id, created_at, type_) VALUES ($1,$2,$3,$4)`
	_, err := r.Db.Exec(query, id, user, now, taskType.String())
	if err != nil {
		return fmt.Errorf("creating initial task record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteEventRepo) ListUserTasks(user u.UserId, p pa.PaginationParams) ([]t.Task, error) {
	q := sq.StatementBuilder.Select(`id,user_id,created_at,type_`).From("tasks")
	q = q.Limit(uint64(p.PageSize)).Offset((uint64(p.Page-1) * uint64(p.PageSize)))
	q = q.Where("user_id=?", user)
	q = q.OrderBy("created_at DESC")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Task{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}
	objects := []t.Task{}
	for _, rec := range records {
		objects = append(objects, t.NewTask(rec.Id, rec.User, rec.Type))
	}

	return objects, nil

}
func (r SQLiteEventRepo) AddEvent(id t.TaskId, event ev.Event) error {
	query := `INSERT INTO events (task_id, time, state, extra, error) VALUES ($1,$2,$3,$4,$5)`
	extraStr, err := serialize(event.Extra)
	if err != nil {
		return fmt.Errorf("creating event record: serializing extra meta-data: %w: %w", err, e.ErrInternal)
	}
	_, err = r.Db.Exec(query, id, event.Time, event.State.String(), extraStr, event.Error)
	if err != nil {
		return fmt.Errorf("creating event record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteEventRepo) GetEvents(id t.TaskId) ([]ev.Event, error) {
	q := sq.StatementBuilder.Select(`task_id,time,state,extra,error`).From("events")
	q = q.Where("task_id=?", id)
	q = q.OrderBy("time DESC")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Event{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}
	events := []ev.Event{}
	for _, rec := range records {
		extra, err := deserialize(rec.Extra)
		if err != nil {
			return nil, fmt.Errorf("deserializing extra meta-data field: %v: %w", err, e.ErrInternal)
		}
		events = append(events, ev.Event{Time: rec.Time, State: rec.State, Extra: extra, Error: rec.Error})
	}

	return events, nil

}
func (r SQLiteEventRepo) ClipNumTasks(user u.UserId, numTasks int) error {
	query := `DELETE FROM events WHERE task_id NOT IN (SELECT id FROM tasks WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2)`
	_, err := r.Db.Exec(query, user, numTasks)
	if err != nil {
		return fmt.Errorf("clipping tasks: %v: %w", err, e.ErrInternal)
	}

	query = `DELETE FROM tasks WHERE id NOT IN (SELECT id FROM tasks WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2)`
	_, err = r.Db.Exec(query, user, numTasks)
	if err != nil {
		return fmt.Errorf("clipping tasks: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func NewSQLiteEventRepo(db adb.Querier) SQLiteEventRepo {
	return SQLiteEventRepo{Db: db}
}
