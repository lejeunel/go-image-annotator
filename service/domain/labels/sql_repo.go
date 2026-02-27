package labels

import (
	e "datahub/errors"
	g "datahub/generic"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SQLLabelRepo struct {
	Db             *sqlx.DB
	ErrorConverter e.DBErrorConverter
}

func NewSQLiteLabelRepo(db *sqlx.DB) *SQLLabelRepo {

	return &SQLLabelRepo{Db: db, ErrorConverter: &e.SQLiteErrorConverter{}}

}

func FromRecord(record LabelRecord) Label {
	return Label{
		Id:          record.Id,
		Name:        record.Name,
		Description: record.Description,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
		ParentId:    record.ParentId,
	}
}

type LabelRecord struct {
	Id          LabelId   `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	ParentId    *LabelId  `db:"parent_id"`
}

func (r *SQLLabelRepo) Create(label *Label) error {

	query := "INSERT INTO labels (id, name, description, created_at, updated_at) VALUES ($1,$2,$3,$4,$5)"
	_, err := r.Db.Exec(query, label.Id, label.Name, label.Description, label.CreatedAt,
		label.UpdatedAt)

	if err != nil {
		err = r.ErrorConverter.Convert(err)
		if errors.Is(err, e.ErrDBUniqueConstraint) {
			err = e.ErrDuplication
		}

		return err
	}

	return nil
}

func (r *SQLLabelRepo) SetParenting(child *Label, parent *Label) error {
	query := "UPDATE labels SET parent_id=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, parent.Id, child.Id)

	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLLabelRepo) Delete(label *Label) error {
	_, err := r.Db.Exec("DELETE FROM labels WHERE id=$1", label.Id)

	err = r.ErrorConverter.Convert(err)

	if errors.Is(err, e.ErrDBForeignKeyConstraint) {
		err = e.ErrDependency
	}
	return err
}

func (r *SQLLabelRepo) Update(label *Label, values Updatables) error {
	query := "UPDATE labels SET description=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, values.Description, label.Id)

	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLLabelRepo) Find(id LabelId) (*Label, error) {
	record := LabelRecord{}
	err := r.Db.Get(&record,
		"SELECT id,parent_id,name,description,created_at,updated_at FROM labels WHERE id=$1", id)

	if err != nil {
		return nil, e.ErrNotFound
	}

	label := FromRecord(record)

	return &label, nil
}

func (r *SQLLabelRepo) FindByName(name string) (*Label, error) {
	query := g.SqlBuilder.Select("id,parent_id,name,description,created_at,updated_at").
		From("labels").
		Where(sq.Eq{"name": name})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", e.ErrDB)
	}

	record := LabelRecord{}
	if err := r.Db.Get(&record, sql, args...); err != nil {
		return nil, fmt.Errorf("slicing labels: %w", r.ErrorConverter.Convert(err))
	}

	label := FromRecord(record)

	return &label, nil
}

func (r *SQLLabelRepo) Count() (int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM labels "
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, e.ErrDB
	}

	return count, nil

}

func (r *SQLLabelRepo) List(ordering g.OrderingArg, pagination g.PaginationParams) ([]Label, *g.PaginationMeta, error) {
	query := g.SqlBuilder.Select("id,parent_id,name,description,created_at,updated_at").From("labels")

	if ordering.Descending == false {
		query = query.OrderBy(ordering.Field + " ASC")
	} else {
		query = query.OrderBy(ordering.Field + " DESC")
	}

	query = query.Limit(uint64(pagination.Limit())).Offset(uint64(pagination.Offset()))
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("building query: %w", e.ErrDB)
	}

	var records []LabelRecord
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, nil, fmt.Errorf("slicing labels: %w", err)
	}

	count, err := r.Count()
	if err != nil {
		return nil, nil, err
	}

	paginationMeta := g.NewPaginationMeta(pagination.Page, count, int64(pagination.PageSize))

	labels := []Label{}
	for _, record := range records {
		label := FromRecord(record)
		labels = append(labels, label)
	}

	return labels, &paginationMeta, nil

}
