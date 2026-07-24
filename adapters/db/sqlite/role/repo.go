package role

import (
	"fmt"

	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	ro "github.com/lejeunel/go-image-annotator/entities/role"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type SQLiteRoleRepo struct {
	Db *sqlx.DB
}

type Row struct {
	Id          ro.RoleId `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}

func (r SQLiteRoleRepo) Create(role ro.Role) error {
	query := `INSERT INTO roles (id, name, description) VALUES ($1,$2,$3)`
	_, err := r.Db.Exec(query, role.Id, role.Name, role.Description)
	if err != nil {
		return fmt.Errorf("creating role: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteRoleRepo) rowToEntity(row Row) ro.Role {
	c := ro.NewRole(row.Id, row.Name,
		ro.WithDescription(row.Description))
	return c

}
func (r SQLiteRoleRepo) Find(name string) (*ro.Role, error) {

	errCtx := fmt.Errorf("fetching role with name %v", name)
	row := Row{}
	err := r.Db.Get(&row,
		`SELECT id,name,description FROM roles WHERE name=$1`, name)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: %w: %w", errCtx, err, e.ErrNotFound)
		default:
			return nil, fmt.Errorf("%w: %w: %w", errCtx, err, e.ErrInternal)
		}
	}

	entity := r.rowToEntity(row)
	return &entity, nil
}
func (r SQLiteRoleRepo) Exists(name string) (*bool, error) {
	var exists bool
	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM roles WHERE name = $1)`, name)
	if err != nil {
		return nil, fmt.Errorf("checking whether role with name %v exists: %v: %w", name, err, e.ErrInternal)
	}

	return &exists, nil
}
func (r SQLiteRoleRepo) Delete(name string) error {
	_, err := r.Db.Exec("DELETE FROM roles WHERE name=$1", name)

	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteRoleRepo) Update(m ro.UpdatableModel) error {
	query := "UPDATE roles SET name=$1,description=$2 WHERE name=$3"
	_, err := r.Db.Exec(query, m.NewName, m.NewDescription, m.Name)

	if err != nil {
		return fmt.Errorf("updating record: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteRoleRepo) IsAssigned(name string) (*bool, error) {
	var count int64

	var isUsed bool
	var query string
	var err error
	query = "SELECT COUNT(*) FROM users_roles WHERE role_id=(SELECT id FROM roles WHERE name=$1)"
	err = r.Db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("checking whether role is assigned to user: %v: %w", err, e.ErrInternal)
	}
	isUsed = count > 0

	return &isUsed, nil
}
func (r SQLiteRoleRepo) List() ([]ro.Role, error) {
	q := sq.StatementBuilder.Select(`id,name,description`).From("roles")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Row{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	objects := []ro.Role{}
	for _, rec := range records {
		e := r.rowToEntity(rec)
		objects = append(objects, e)
	}

	return objects, nil
}

func NewSQLiteRoleRepo(db *sqlx.DB) SQLiteRoleRepo {
	return SQLiteRoleRepo{Db: db}
}

func NewTestSQLiteRoleRepo() SQLiteRoleRepo {
	return NewSQLiteRoleRepo(s.NewSQLiteDB(":memory:"))
}
