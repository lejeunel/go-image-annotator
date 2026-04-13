package image

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	ist "github.com/lejeunel/go-image-annotator-v2/application/image-store"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	s "github.com/lejeunel/go-image-annotator-v2/infra/db"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

type SQLiteImageRepo struct {
	Db *sqlx.DB
}

type Row struct {
	ImageId      im.ImageId       `db:"image_id"`
	CollectionId clc.CollectionId `db:"collection_id"`
}

func (r *SQLiteImageRepo) AddToCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	query := "INSERT INTO images_collections (image_id, collection_id) VALUES ($1,$2)"
	_, err := r.Db.Exec(query, imageId.String(), collectionId.String())
	if err != nil {
		return fmt.Errorf("inserting record into image to collection junction table: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r *SQLiteImageRepo) Count(f ist.CountingParams) (*int64, error) {
	var count int64

	var query string
	var err error
	if f.Collection != nil {
		query = "SELECT COUNT(*) FROM images_collections WHERE collection_id=(SELECT id FROM collections WHERE name=$1)"
		err = r.Db.QueryRow(query, f.Collection).Scan(&count)
	} else {
		query = "SELECT COUNT(*) FROM images"
		err = r.Db.QueryRow(query).Scan(&count)
	}
	if err != nil {
		return nil, fmt.Errorf("counting image records: %v: %w", err, e.ErrInternal)
	}
	return &count, nil

}

func (r *SQLiteImageRepo) List(f ist.FilteringParams) (*[]im.BaseImage, error) {

	q := sq.StatementBuilder.Select("image_id,collection_id").From("images_collections")
	q = q.Limit(uint64(f.PageSize)).Offset((uint64(f.Page-1) * uint64(f.PageSize)))

	if f.Collection != nil {
		q = q.Where(fmt.Sprintf("collection_id=(SELECT id FROM collections WHERE name='%v')", *f.Collection))
	}

	q = q.OrderBy("image_id")
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Row{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	objects := []im.BaseImage{}
	for _, rec := range records {
		var collectionName string
		q := "SELECT name FROM collections WHERE id=$1"
		err := r.Db.QueryRow(q, rec.CollectionId.String()).Scan(&collectionName)
		if err != nil {
			return nil, fmt.Errorf("fetching collection name from id %v: %v: %w", rec.CollectionId, err, e.ErrInternal)
		}
		objects = append(objects, im.BaseImage{ImageId: rec.ImageId, Collection: collectionName})
	}

	return &objects, nil
}
func (r *SQLiteImageRepo) ImageExistsInCollection(imageId im.ImageId, collectionId clc.CollectionId) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM images_collections WHERE image_id=$1 AND collection_id=$2"
	err := r.Db.QueryRow(query, imageId.String(), collectionId.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking image to collection junction records: %v: %w", err, e.ErrInternal)
	}

	return count > 0, nil
}
func (r *SQLiteImageRepo) ImageExists(imageId im.ImageId) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM images WHERE id=$1"
	err := r.Db.QueryRow(query, imageId.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking that image exists: %v: %w", err, e.ErrInternal)
	}

	return count > 0, nil
}
func (r *SQLiteImageRepo) MIMEType(imageId im.ImageId) (*string, error) {
	errCtx := "finding image MIMEType"
	var mimetype string
	err := r.Db.Get(&mimetype, "SELECT mimetype FROM images WHERE id = $1", imageId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrNotFound)
		default:
			return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
		}
	}
	return &mimetype, nil
}

func (r *SQLiteImageRepo) AddImage(imageId im.ImageId, hash []byte, format string) error {
	query := "INSERT INTO images (id, hash, mimetype) VALUES ($1,$2,$3)"
	_, err := r.Db.Exec(query, imageId.String(), hex.EncodeToString(hash), format)
	if err != nil {
		return fmt.Errorf("inserting image record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteImageRepo) FindImageIdByHash(hash []byte) (*im.ImageId, error) {
	errCtx := "finding image record by hash"
	var imageId im.ImageId
	err := r.Db.Get(&imageId, "SELECT id FROM images WHERE hash = $1", hex.EncodeToString(hash))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrNotFound)
		}
		return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
	}
	return &imageId, nil
}
func (r *SQLiteImageRepo) Delete(id im.ImageId) error {
	_, err := r.Db.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleting image record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteImageRepo) RemoveImageFromCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	_, err := r.Db.Exec("DELETE FROM images_collections WHERE image_id = $1 AND collection_id = $2",
		imageId, collectionId)
	if err != nil {
		return fmt.Errorf("removing image from image to collection junction table: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func NewSQLiteImageRepo(db *sqlx.DB) *SQLiteImageRepo {
	return &SQLiteImageRepo{Db: db}
}

func NewTestSQLiteImageRepo() *SQLiteImageRepo {
	return NewSQLiteImageRepo(s.NewSQLiteDB(":memory:"))
}
