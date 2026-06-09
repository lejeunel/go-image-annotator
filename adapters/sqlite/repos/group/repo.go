package group

import (
	"fmt"

	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	s "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/group/list"
	"github.com/lejeunel/go-image-annotator/use-cases/group/update"
)

type SQLiteGroupRepo struct {
	Db *sqlx.DB
}

type Row struct {
	Id          g.GroupId `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}

func (r *SQLiteGroupRepo) Create(grp g.Group) error {
	query := `INSERT INTO groups (id, name, description) VALUES ($1,$2,$3)`
	_, err := r.Db.Exec(query, grp.Id.String(), grp.Name, grp.Description)
	if err != nil {
		return fmt.Errorf("creating record: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r *SQLiteGroupRepo) rowToEntity(row Row) g.Group {
	c := g.NewGroup(row.Id, row.Name,
		g.WithDescription(row.Description))
	return c

}
func (r *SQLiteGroupRepo) Find(name string) (*g.Group, error) {

	row := Row{}
	err := r.Db.Get(&row,
		`SELECT id,name,description FROM groups WHERE name=$1`, name)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrNotFound
		default:
			return nil, fmt.Errorf("fetching record by name: %v: %w", err, e.ErrInternal)
		}
	}

	entity := r.rowToEntity(row)
	return &entity, nil
}
func (r *SQLiteGroupRepo) Exists(name string) (*bool, error) {
	var exists bool

	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM groups WHERE name = $1)`, name)
	if err != nil {
		return nil, fmt.Errorf("checking whether record exists: %v: %w", err, e.ErrInternal)
	}

	return &exists, nil
}
func (r *SQLiteGroupRepo) Delete(name string) error {
	_, err := r.Db.Exec("DELETE FROM groups WHERE name=$1", name)

	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteGroupRepo) Update(m update.Model) error {
	query := "UPDATE groups SET name=$1,description=$2 WHERE name=$3"
	_, err := r.Db.Exec(query, m.NewName, m.NewDescription, m.Name)

	if err != nil {
		return fmt.Errorf("updating record: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r *SQLiteGroupRepo) IsUsed(name string) (*bool, error) {
	var count int64

	var isUsed bool
	var query string
	var err error
	query = "SELECT COUNT(*) FROM users_groups WHERE group_id=(SELECT id FROM groups WHERE name=$1)"
	err = r.Db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("checking whether group is assigned to user: %v: %w", err, e.ErrInternal)
	}
	isUsed = count > 0

	query = "SELECT COUNT(*) FROM collections WHERE group_id=(SELECT id FROM groups WHERE name=$1)"
	err = r.Db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("checking whether group is assigned to collection: %v: %w", err, e.ErrInternal)
	}
	isUsed = isUsed || count > 0

	return &isUsed, nil
}
func (r *SQLiteGroupRepo) Count() (*int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM groups"
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("counting records: %v: %w", err, e.ErrInternal)
	}

	return &count, nil
}
func (r *SQLiteGroupRepo) GroupOfCollection(name string) (*string, error) {
	var group string

	err := r.Db.Get(&group, `SELECT name FROM groups WHERE id=(SELECT group_id FROM collections WHERE name=$1)`, name)
	if err != nil {
		return nil, fmt.Errorf("checking whether record exists: %v: %w", err, e.ErrInternal)
	}

	return &group, nil

}
func (r *SQLiteGroupRepo) List(m list.Request) ([]*g.Group, error) {
	q := sq.StatementBuilder.Select(`id,name,description`).From("groups")
	q = q.Limit(uint64(m.PageSize)).Offset((uint64(m.Page-1) * uint64(m.PageSize)))
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Row{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	objects := []*g.Group{}
	for _, rec := range records {
		e := r.rowToEntity(rec)
		objects = append(objects, &e)
	}

	return objects, nil
}

func NewSQLiteGroupRepo(db *sqlx.DB) *SQLiteGroupRepo {
	return &SQLiteGroupRepo{Db: db}
}

func NewTestSQLiteGroupRepo() *SQLiteGroupRepo {
	return NewSQLiteGroupRepo(s.NewSQLiteDB(":memory:"))
}
