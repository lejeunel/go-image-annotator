package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	"time"
)

type SQLImageRepo struct {
	Db *sqlx.DB
}

func NewSQLImageRepo(db *sqlx.DB) *SQLImageRepo {

	return &SQLImageRepo{Db: db}

}

func (r SQLImageRepo) Create(ctx context.Context, image *m.Image) (*m.Image, error) {
	now := time.Now().String()
	query := "INSERT INTO images (id, uri, created_at, updated_at, sha256, width, height, mimetype) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := r.Db.Exec(query, image.Id, image.Uri, now,
		now, image.SHA256, image.Width, image.Height, image.MIMEType)

	if err != nil {
		return nil, err
	}

	return image, nil
}

func (r SQLImageRepo) GetOne(ctx context.Context, id string) (*m.Image, error) {
	image := m.Image{}
	err := r.Db.Get(&image, "SELECT id,uri,created_at,updated_at,sha256,width,height,mimetype FROM images WHERE id=?", id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "image", Criteria: "id", Value: id, Err: err}
	}

	return &image, nil
}

func (r *SQLImageRepo) Nums() (int64, error) {

	return 0, nil

}

func (r *SQLImageRepo) Slice(offset, length int, data interface{}) error {

	return nil
}
