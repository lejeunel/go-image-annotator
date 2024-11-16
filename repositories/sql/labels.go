package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	"time"
)

type SQLLabelRepo struct {
	Db *sqlx.DB
}

func NewSQLLabelRepo(db *sqlx.DB) *SQLLabelRepo {

	return &SQLLabelRepo{Db: db}

}

func (r SQLLabelRepo) Create(ctx context.Context, label *m.Label) (*m.Label, error) {
	now := time.Now().String()
	query := "INSERT INTO labels (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := r.Db.Exec(query, label.Id, label.Name, label.Description, now,
		now)

	if err != nil {
		return nil, err
	}

	return label, nil
}

func (r SQLLabelRepo) GetOne(ctx context.Context, id string) (*m.Label, error) {
	label := m.Label{}
	err := r.Db.Get(&label, "SELECT id,name,description FROM labels WHERE id=?", id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "label", Criteria: "id", Value: id, Err: err}
	}

	return &label, nil
}

func (r *SQLLabelRepo) Nums() (int64, error) {

	return 0, nil

}

func (r *SQLLabelRepo) Slice(offset, length int, data interface{}) error {

	return nil
}
