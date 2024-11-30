package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	"time"
)

type SQLCollectionRepo struct {
	Db *sqlx.DB
}

func NewSQLSetRepo(db *sqlx.DB) *SQLCollectionRepo {

	return &SQLCollectionRepo{Db: db}

}

func (r *SQLCollectionRepo) Create(ctx context.Context, set *m.Collection) (*m.Collection, error) {
	now := time.Now().String()
	query := "INSERT INTO collections (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)"
	_, err := r.Db.Exec(query, set.Id, set.Name, now, now)

	if err != nil {
		return nil, err
	}

	return set, nil
}

func (r *SQLCollectionRepo) Get(ctx context.Context, id string) (*m.Collection, error) {
	set := m.Collection{}
	err := r.Db.Get(&set, "SELECT id,name,created_at,updated_at FROM collections WHERE id=?", id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "collection", Criteria: "id", Value: id, Err: err}
	}

	return &set, nil
}

func (r *SQLCollectionRepo) Delete(ctx context.Context, collection *m.Collection) error {
	id := collection.Id.String()
	_, err := r.Db.Exec("DELETE FROM collections WHERE id=?", id)

	if err != nil {
		return &e.ErrNotFound{Entity: "collection", Criteria: "id", Value: id, Err: err}
	}

	return nil
}

func (r *SQLCollectionRepo) AssignImageToCollection(ctx context.Context, image *m.Image, set *m.Collection) error {
	now := time.Now().String()
	query := "INSERT INTO image_collection_assoc (id, image_id, collection_id, created_at) VALUES (?, ?, ?, ?)"
	_, err := r.Db.Exec(query, uuid.New(), image.Id, set.Id, now)

	if err != nil {
		return err
	}

	return nil

}

func (r *SQLCollectionRepo) Nums() (int64, error) {

	return 0, nil

}

func (r *SQLCollectionRepo) Slice(offset, length int, data interface{}) error {

	return nil
}
