package profile

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	g "github.com/lejeunel/go-image-annotator/entities/group"
	pr "github.com/lejeunel/go-image-annotator/entities/profile"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ProfileRepo struct {
	Db adb.Querier
}

type Row struct {
	Id          pr.ProfileId `db:"id"`
	Name        string       `db:"name"`
	Description string       `db:"description"`
	GroupId     *g.GroupId   `db:"group_id"`
	GroupName   *string      `db:"group_name"`
}

func (r ProfileRepo) Create(p pr.Profile) error {
	var err error
	if p.Group != nil {
		query := `INSERT INTO profiles (id, name, description, group_id) VALUES ($1,$2,$3,(SELECT id FROM groups WHERE name=$4))`
		_, err = r.Db.Exec(query, p.Id, p.Name, p.Description, *p.Group)
	} else {
		query := `INSERT INTO profiles (id, name, description) VALUES ($1,$2,$3)`
		_, err = r.Db.Exec(query, p.Id, p.Name, p.Description)
	}
	if err != nil {
		return fmt.Errorf("creating profile record: %v: %w", err, e.ErrInternal)
	}
	if err := r.addLabels(p.Name, p.Labels); err != nil {
		return fmt.Errorf("adding labels to profile: %w", err)
	}
	return nil
}

func (r ProfileRepo) Find(name string) (*pr.Profile, error) {
	row := Row{}
	errCtx := fmt.Errorf("fetching profile record")
	err := r.Db.Get(&row,
		`
		SELECT c.id,c.name,c.description,c.group_id,g.name AS group_name
		FROM profiles AS c
		LEFT JOIN groups g ON g.id = c.group_id
		WHERE c.name=$1`, name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: %w", err, e.ErrNotFound)
		default:
			return nil, fmt.Errorf("%w: %v: %w", errCtx, err, e.ErrInternal)
		}
	}

	entity, err := r.build(row)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r ProfileRepo) Exists(name string) (*bool, error) {
	var exists bool

	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM profiles WHERE name = $1)`, name)
	if err != nil {
		return &exists, fmt.Errorf("checking whether record exists: %v: %w", err, e.ErrInternal)
	}
	return &exists, nil
}

func (r ProfileRepo) Delete(name string) error {
	_, err := r.Db.Exec("DELETE FROM profiles WHERE name=$1", name)
	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r ProfileRepo) Update(m pr.UpdateModel) error {
	errCtx := fmt.Errorf("updating profile record")
	var err error
	if m.NewGroup != nil {
		query := "UPDATE profiles SET name=$1,description=$2,group_id=(SELECT id FROM groups WHERE name=$3) WHERE name=$4"
		_, err = r.Db.Exec(query, m.NewName, m.NewDescription, *m.NewGroup, m.Name)
	} else {
		query := "UPDATE profiles SET name=$1,description=$2,group_id=NULL WHERE name=$3"
		_, err = r.Db.Exec(query, m.NewName, m.NewDescription, m.Name)
	}

	if err != nil {
		return fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
	}

	if err := r.removeLabels(m.NewName); err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}
	if err := r.addLabels(m.NewName, m.NewLabels); err != nil {
		return fmt.Errorf("%w: %w", errCtx, err)
	}

	return nil
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

func (r ProfileRepo) List(m pa.PaginationParams) ([]pr.Profile, error) {
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

	objects := []pr.Profile{}
	for _, rec := range records {
		obj, err := r.build(rec)
		if err != nil {
			return nil, err
		}
		objects = append(objects, *obj)
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
	var isUsed bool
	var count int64
	query := "SELECT COUNT(*) FROM collections WHERE profile_id=(SELECT id FROM profiles where name=$1)"
	err := r.Db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf(
			"checking whether profile %v is used: %v: %w",
			name,
			err,
			e.ErrInternal,
		)
	}
	if count == 0 {
		isUsed = false
	} else {
		isUsed = true
	}
	return &isUsed, nil
}

func (r ProfileRepo) removeLabels(name pr.ProfileName) error {
	_, err := r.Db.Exec(
		"DELETE FROM profiles_labels WHERE profile_id=(SELECT id FROM profiles WHERE name=$1)",
		name,
	)
	if err != nil {
		return fmt.Errorf("deleting labels from profile %v: %v: %w", name, err, e.ErrInternal)
	}
	return nil
}

func (r ProfileRepo) addLabels(profile pr.ProfileName, labels []string) error {
	for _, label := range labels {
		_, err := r.Db.Exec(
			`INSERT INTO profiles_labels (profile_id,label_id) VALUES ((SELECT id FROM profiles WHERE name=$1), (SELECT id FROM labels WHERE name=$2))`,
			profile,
			label,
		)
		if err != nil {
			return fmt.Errorf(
				"creating profile record: adding label %v: %v: %w",
				label,
				err,
				e.ErrInternal,
			)
		}
	}
	return nil
}

func (r ProfileRepo) build(row Row) (*pr.Profile, error) {
	var labels []string
	if err := r.Db.Select(
		&labels,
		`SELECT name FROM labels WHERE id IN (SELECT label_id FROM profiles_labels WHERE profile_id=$1)`,
		row.Id,
	); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}
	p := pr.NewProfile(row.Id, row.Name,
		pr.WithDescription(row.Description), pr.WithLabels(labels))
	if row.GroupName != nil {
		p.Group = row.GroupName
	}
	return &p, nil
}

func NewProfileRepo(db adb.Querier) ProfileRepo {
	return ProfileRepo{Db: db}
}
