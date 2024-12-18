package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	pag "github.com/vcraescu/go-paginator/v2"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	"time"
)

type SQLCollectionRepo struct {
	Db *sqlx.DB
}

type PaginableSQLCollectionRepo struct {
	Repo *SQLCollectionRepo
}

func NewSQLCollectionRepo(db *sqlx.DB) *SQLCollectionRepo {

	return &SQLCollectionRepo{Db: db}

}

func (r *SQLCollectionRepo) Create(ctx context.Context, collection *m.Collection) (*m.Collection, error) {
	now := time.Now()

	query := "INSERT INTO collections (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)"
	_, err := r.Db.Exec(query, collection.Id, collection.Name, now, now)

	if err != nil {
		return nil, err
	}

	return collection, nil
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

func (r *SQLCollectionRepo) AssignImageToCollection(ctx context.Context, image *m.Image, collection *m.Collection) error {
	now := time.Now()
	query := "INSERT INTO image_collection_assoc (id, image_id, collection_id, created_at) VALUES (?, ?, ?, ?)"
	_, err := r.Db.Exec(query, uuid.New(), image.Id, collection.Id, now)

	if err != nil {
		return err
	}

	return nil

}

func (r *SQLCollectionRepo) ImageIsInCollection(ctx context.Context, image *m.Image, collection *m.Collection) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM image_collection_assoc WHERE image_id=? AND collection_id=?"
	err := r.Db.QueryRow(query, image.Id, collection.Id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *SQLCollectionRepo) RemoveImage(ctx context.Context, image *m.Image, collection *m.Collection) error {
	collection_id := collection.Id.String()
	image_id := image.Id.String()
	_, err := r.Db.Exec("DELETE FROM image_collection_assoc WHERE image_id=? AND collection_id=?",
		image_id, collection_id)

	if err != nil {
		return err
	}

	return nil

}

func (r *SQLCollectionRepo) Paginate(pageSize int) pag.Paginator {
	paginable := &PaginableSQLCollectionRepo{Repo: r}
	return pag.New(paginable, pageSize)

}

func (r *PaginableSQLCollectionRepo) Nums() (int64, error) {

	var count int64

	query := "SELECT COUNT(*) FROM images"
	err := r.Repo.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func (r *PaginableSQLCollectionRepo) Slice(offset, length int, data interface{}) error {

	query := "SELECT id,name,created_at,updated_at FROM collections " + " LIMIT $1 OFFSET $2"
	rows, err := r.Repo.Db.Queryx(query, length, offset)

	if err != nil {
		return err
	}

	s := data.(*[]m.Collection)

	for rows.Next() {
		var c m.Collection
		err := rows.StructScan(&c)

		if err != nil {
			return err
		}

		*s = append(*s, c)
	}

	return nil
}
