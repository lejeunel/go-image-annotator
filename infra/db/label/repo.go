package label

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	s "github.com/lejeunel/go-image-annotator/infra/db"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/label/list"
	"github.com/lejeunel/go-image-annotator/use-cases/label/update"
)

type SQLiteLabelRepo struct {
	Db *sqlx.DB
}

type LabelRecord struct {
	Id          lbl.LabelId `db:"id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
}

func (r *SQLiteLabelRepo) Create(l lbl.Label) error {
	query := "INSERT INTO labels (id, name, description) VALUES ($1,$2,$3)"
	_, err := r.Db.Exec(query, l.Id.String(), l.Name, l.Description)
	if err != nil {
		return fmt.Errorf("inserting record: %v: %w", err, e.ErrInternal)
	}

	return nil

}

func (r *SQLiteLabelRepo) FindLabelByName(name string) (*lbl.Label, error) {
	record := LabelRecord{}
	err := r.Db.Get(&record,
		"SELECT id,name,description FROM labels WHERE name=$1", name)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrNotFound
		default:
			return nil, fmt.Errorf("finding record by label: %v: %w", err, e.ErrInternal)
		}

	}

	return lbl.NewLabel(record.Id, record.Name, lbl.WithDescription(record.Description)), nil
}

func (r *SQLiteLabelRepo) Delete(name string) error {
	_, err := r.Db.Exec("DELETE FROM labels WHERE name=$1", name)

	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteLabelRepo) Count() (int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM labels"
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting records: %v: %w", err, e.ErrInternal)
	}

	return count, nil
}

func (r *SQLiteLabelRepo) List(m list.Request) ([]*lbl.Label, error) {
	q := sq.StatementBuilder.Select("id,name,description").From("labels")
	q = q.Limit(uint64(m.PageSize)).Offset((uint64(m.Page-1) * uint64(m.PageSize)))
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []LabelRecord{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	objects := []*lbl.Label{}
	for _, r := range records {
		objects = append(objects, &lbl.Label{Id: r.Id, Name: r.Name, Description: r.Description})
	}

	return objects, nil
}
func (r *SQLiteLabelRepo) FetchAll() ([]string, error) {
	q := sq.StatementBuilder.Select("name").From("labels")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	names := []string{}
	if err := r.Db.Select(&names, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	return names, nil
}

func (r *SQLiteLabelRepo) Update(m update.Model) error {
	query := "UPDATE labels SET name=$1,description=$2 WHERE name=$3"
	_, err := r.Db.Exec(query, m.NewName, m.NewDescription, m.Name)

	if err != nil {
		return fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return nil
}

func (r *SQLiteLabelRepo) Exists(name string) (bool, error) {
	var exists bool

	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM labels WHERE name = $1)`, name)
	if err != nil {
		return false, fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return exists, nil
}

func (r *SQLiteLabelRepo) IsUsed(name string) (*bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM annotations WHERE label_id=(SELECT id FROM labels WHERE name=$1)"
	err := r.Db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	isUsed := count > 0
	return &isUsed, nil

}

func NewSQLiteLabelRepo(db *sqlx.DB) *SQLiteLabelRepo {
	return &SQLiteLabelRepo{Db: db}
}

func NewTestSQLiteLabelRepo() *SQLiteLabelRepo {
	return NewSQLiteLabelRepo(s.NewSQLiteDB(":memory:"))
}
