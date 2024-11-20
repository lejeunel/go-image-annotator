package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	pag "github.com/vcraescu/go-paginator/v2"
	e "go-image-annotator/errors"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	"strings"
	"time"
)

type SQLImageRepo struct {
	Db *sqlx.DB
}

type PaginableSQLImageRepo struct {
	Repo    *SQLImageRepo
	Filters *g.ImageFilterArgs
}

func NewSQLImageRepo(db *sqlx.DB) *SQLImageRepo {

	return &SQLImageRepo{Db: db}

}

func (r *SQLImageRepo) Create(ctx context.Context, image *m.Image) (*m.Image, error) {
	now := time.Now().String()
	query := "INSERT INTO images (id, uri, created_at, updated_at, sha256, width, height, mimetype) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := r.Db.Exec(query, image.Id, image.Uri, now,
		now, image.SHA256, image.Width, image.Height, image.MIMEType)

	if err != nil {
		return nil, err
	}

	return image, nil
}

func (r *SQLImageRepo) GetOne(ctx context.Context, id string) (*m.Image, error) {
	image := m.Image{}
	err := r.Db.Get(&image, "SELECT id,uri,created_at,updated_at,sha256,width,height,mimetype FROM images WHERE id=?", id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "image", Criteria: "id", Value: id, Err: err}
	}

	return &image, nil
}

func (r *SQLImageRepo) Delete(ctx context.Context, image *m.Image) error {
	_, err_image := r.Db.Exec("DELETE FROM images WHERE id=?", image.Id.String())
	_, err_assoc := r.Db.Exec("DELETE FROM image_label_assoc WHERE image_id=?", image.Id.String())
	return errors.Join(err_image, err_assoc)
}

func (r *SQLImageRepo) Paginate(pageSize int, filters *g.ImageFilterArgs) pag.Paginator {
	paginable := &PaginableSQLImageRepo{Repo: r, Filters: filters}
	return pag.New(paginable, pageSize)

}

func (r *PaginableSQLImageRepo) buildFilteringWhereClause() (string, error) {
	var parts []string
	if r.Filters.SetId != "" {
		parts = append(parts, fmt.Sprintf("id in (SELECT image_id FROM image_set_assoc WHERE set_id=\"%v\")", r.Filters.SetId))
	}

	if r.Filters.SetName != "" {
		set := m.Set{}
		err := r.Repo.Db.Get(&set, "SELECT id FROM imagesets WHERE name=?", r.Filters.SetName)
		if err != nil {
			return "", err
		}

		parts = append(parts, fmt.Sprintf("id in (SELECT image_id FROM image_set_assoc WHERE set_id=\"%v\")", set.Id.String()))
	}

	clause := ""
	if len(parts) > 0 {
		clause = "WHERE " + strings.Join(parts, " ")
	}

	return clause, nil

}

func (r *PaginableSQLImageRepo) Nums() (int64, error) {
	var count int64

	filteringWhereClause, err := r.buildFilteringWhereClause()
	if err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM images " + filteringWhereClause
	err = r.Repo.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func (r *PaginableSQLImageRepo) Slice(offset, length int, data interface{}) error {
	filteringWhereClause, err := r.buildFilteringWhereClause()
	if err != nil {
		return err
	}

	baseQuery := "SELECT id,uri,created_at,updated_at,sha256,width,height,mimetype FROM images "
	query := baseQuery + filteringWhereClause + " LIMIT $1 OFFSET $2"
	rows, err := r.Repo.Db.Queryx(query, length, offset)

	if err != nil {
		return err
	}

	s := data.(*[]m.Image)

	for rows.Next() {
		var b m.Image
		err := rows.StructScan(&b)

		if err != nil {
			return err
		}

		*s = append(*s, b)
	}

	return nil
}
