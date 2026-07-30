package collection

import (
	"fmt"

	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type SQLiteCollectionRepo struct {
	Db adb.Querier
}

type Row struct {
	Id          clc.CollectionId `db:"id"`
	Name        string           `db:"name"`
	Description string           `db:"description"`
	CreatedAt   sql.NullTime     `db:"created_at"`
	GroupId     *g.GroupId       `db:"group_id"`
	GroupName   *string          `db:"group_name"`
}

func (r SQLiteCollectionRepo) Create(c clc.Collection) error {
	var err error
	if c.Group != nil {
		query := `INSERT INTO collections (id, name, description, created_at, group_id) VALUES ($1,$2,$3,$4,(SELECT id FROM groups WHERE name=$5))`
		_, err = r.Db.Exec(query, c.Id.String(), c.Name, c.Description, c.CreatedAt, *c.Group)
	} else {
		query := `INSERT INTO collections (id, name, description, created_at) VALUES ($1,$2,$3,$4)`
		_, err = r.Db.Exec(query, c.Id.String(), c.Name, c.Description, c.CreatedAt)
	}
	if err != nil {
		return fmt.Errorf("creating record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteCollectionRepo) rowToEntity(row Row) clc.Collection {
	c := clc.NewCollection(row.Id, row.Name,
		clc.WithDescription(row.Description))
	if row.CreatedAt.Valid {
		c.CreatedAt = row.CreatedAt.Time
	}
	if row.GroupName != nil {
		c.Group = row.GroupName
	}
	return c

}
func (r SQLiteCollectionRepo) Find(name string) (*clc.Collection, error) {

	row := Row{}
	err := r.Db.Get(&row,
		`
		SELECT c.id,c.name,c.description,c.created_at,c.group_id,g.name AS group_name
		FROM collections AS c
		LEFT JOIN groups g ON g.id = c.group_id
		WHERE c.name=$1`, name)

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
func (r SQLiteCollectionRepo) Exists(name string) (bool, error) {
	var exists bool

	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM collections WHERE name = $1)`, name)
	if err != nil {
		return false, fmt.Errorf("checking whether record exists: %v: %w", err, e.ErrInternal)
	}

	return exists, nil
}
func (r SQLiteCollectionRepo) Delete(name string) error {
	_, err := r.Db.Exec("DELETE FROM collections WHERE name=$1", name)

	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteCollectionRepo) Update(m clc.UpdateModel) error {
	var err error
	if m.NewGroup != nil {
		fmt.Println("repo updating with group", *m.NewGroup)
		query := "UPDATE collections SET name=$1,description=$2,group_id=(SELECT id FROM groups WHERE name=$3) WHERE name=$4"
		_, err = r.Db.Exec(query, m.NewName, m.NewDescription, *m.NewGroup, m.Name)
	} else {
		query := "UPDATE collections SET name=$1,description=$2,group_id=NULL WHERE name=$3"
		_, err = r.Db.Exec(query, m.NewName, m.NewDescription, m.Name)
	}

	if err != nil {
		return fmt.Errorf("updating record: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteCollectionRepo) IsPopulated(name string) (*bool, error) {
	var count int64

	var query string
	var err error
	query = "SELECT COUNT(*) FROM images_collections WHERE collection_id=(SELECT id FROM collections WHERE name=$1)"
	err = r.Db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("checking whether collection is populated: %v: %w", err, e.ErrInternal)
	}
	isPopulated := count > 0
	return &isPopulated, nil
}
func (r SQLiteCollectionRepo) Count() (*int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM collections"
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("counting records: %v: %w", err, e.ErrInternal)
	}

	return &count, nil
}
func (r SQLiteCollectionRepo) List(m pa.PaginationParams) ([]*clc.Collection, error) {
	q := sq.StatementBuilder.Select(`c.id,c.name,c.description,c.created_at,c.group_id,g.name AS group_name`).From("collections AS c")
	q = q.LeftJoin("groups g ON g.id=c.group_id")
	q = q.Limit(uint64(m.PageSize)).Offset((uint64(m.Page-1) * uint64(m.PageSize)))
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Row{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	objects := []*clc.Collection{}
	for _, rec := range records {
		e := r.rowToEntity(rec)
		objects = append(objects, &e)
	}

	return objects, nil
}
func (r SQLiteCollectionRepo) GetGroup(name string) (*string, error) {
	var group string
	errCtx := fmt.Errorf("retrieving group of collection with name %v", name)

	err := r.Db.Get(&group, `SELECT name FROM groups WHERE id=(SELECT group_id FROM collections WHERE name=$1)`, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", errCtx, e.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %w: %w", errCtx, err, e.ErrInternal)
	}

	return &group, nil

}

func NewSQLiteCollectionRepo(db adb.Querier) SQLiteCollectionRepo {
	return SQLiteCollectionRepo{Db: db}
}
