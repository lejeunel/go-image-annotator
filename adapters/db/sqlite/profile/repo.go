package profile

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	clc "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ProfileRepo struct {
	Db adb.Querier
}

type Row struct {
	Id          clc.ProfileId `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	GroupId     *g.GroupId    `db:"group_id"`
	GroupName   *string       `db:"group_name"`
}

func (r ProfileRepo) Create(c clc.Profile) error {
	var err error
	if c.Group != nil {
		query := `INSERT INTO profiles (id, name, description, group_id) VALUES ($1,$2,$3,(SELECT id FROM groups WHERE name=$4))`
		_, err = r.Db.Exec(query, c.Id.String(), c.Name, c.Description, *c.Group)
	} else {
		query := `INSERT INTO profiles (id, name, description) VALUES ($1,$2,$3)`
		_, err = r.Db.Exec(query, c.Id.String(), c.Name, c.Description)
	}
	if err != nil {
		return fmt.Errorf("creating record: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r ProfileRepo) build(row Row) clc.Profile {
	c := clc.NewProfile(row.Id, row.Name,
		clc.WithDescription(row.Description))
	if row.GroupName != nil {
		c.Group = row.GroupName
	}
	return c
}

func (r ProfileRepo) Find(name string) (*clc.Profile, error) {
	row := Row{}
	err := r.Db.Get(&row,
		`
		SELECT c.id,c.name,c.description,c.group_id,g.name AS group_name
		FROM profiles AS c
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

	entity := r.build(row)

	return &entity, nil
}

func (r ProfileRepo) Exists(name string) (bool, error) {
	var exists bool

	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM profiles WHERE name = $1)`, name)
	if err != nil {
		return false, fmt.Errorf("checking whether record exists: %v: %w", err, e.ErrInternal)
	}

	return exists, nil
}

func (r ProfileRepo) Delete(name string) error {
	_, err := r.Db.Exec("DELETE FROM profiles WHERE name=$1", name)
	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r ProfileRepo) Update(m clc.UpdateModel) error {
	var err error
	if m.NewGroup != nil {
		query := "UPDATE profiles SET name=$1,description=$2,group_id=(SELECT id FROM groups WHERE name=$3) WHERE name=$4"
		_, err = r.Db.Exec(query, m.NewName, m.NewDescription, *m.NewGroup, m.Name)
	} else {
		query := "UPDATE profiles SET name=$1,description=$2,group_id=NULL WHERE name=$3"
		_, err = r.Db.Exec(query, m.NewName, m.NewDescription, m.Name)
	}

	if err != nil {
		return fmt.Errorf("updating record: %v: %w", err, e.ErrInternal)
	}

	return nil
}

func (r ProfileRepo) IsPopulated(name string) (*bool, error) {
	var count int64

	var query string
	var err error
	query = "SELECT COUNT(*) FROM images_profiles WHERE profile_id=(SELECT id FROM profiles WHERE name=$1)"
	err = r.Db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf(
			"checking whether profile is populated: %v: %w",
			err,
			e.ErrInternal,
		)
	}
	isPopulated := count > 0
	return &isPopulated, nil
}

func (r ProfileRepo) Count() (*int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM profiles"
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("counting records: %v: %w", err, e.ErrInternal)
	}

	return &count, nil
}

func (r ProfileRepo) List(m pa.PaginationParams) ([]*clc.Profile, error) {
	q := sq.StatementBuilder.Select(`c.id,c.name,c.description,c.group_id,g.name AS group_name`).
		From("profiles AS c")
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

	objects := []*clc.Profile{}
	for _, rec := range records {
		e := r.build(rec)
		objects = append(objects, &e)
	}

	return objects, nil
}

func (r ProfileRepo) GetGroup(name string) (*string, error) {
	var group string
	errCtx := fmt.Errorf("retrieving group of profile with name %v", name)

	err := r.Db.Get(
		&group,
		`SELECT name FROM groups WHERE id=(SELECT group_id FROM profiles WHERE name=$1)`,
		name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", errCtx, e.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %w: %w", errCtx, err, e.ErrInternal)
	}

	return &group, nil
}

func (r ProfileRepo) IsUsed(name string) (*bool, error) {
	res := true
	return &res, nil
}

func NewProfileRepo(db adb.Querier) ProfileRepo {
	return ProfileRepo{Db: db}
}
