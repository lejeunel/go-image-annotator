package image

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"iter"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	s "github.com/lejeunel/go-image-annotator/adapters/db/sqlite"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	"time"
)

type SQLiteImageRepo struct {
	Db *sqlx.DB
}

type ListRow struct {
	ImageId      im.ImageId       `db:"image_id"`
	CollectionId clc.CollectionId `db:"collection_id"`
	Name         string           `db:"name"`
	IngestTime   time.Time        `db:"ingested_at"`
}

type SpecsRow struct {
	MIMEType   string    `db:"mimetype"`
	Width      int       `db:"width"`
	Height     int       `db:"height"`
	IngestedAt time.Time `db:"ingested_at"`
}

func (r SQLiteImageRepo) AddToCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	query := "INSERT INTO images_collections (image_id, collection_id) VALUES ($1,$2)"
	_, err := r.Db.Exec(query, imageId.String(), collectionId.String())
	if err != nil {
		return fmt.Errorf("inserting record into image to collection junction table: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteImageRepo) Count(f im.CountingParams) (*int64, error) {
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
func (r SQLiteImageRepo) Slice(f im.FilteringParams, p pa.PaginationParams, o im.OrderingParams) ([]im.BaseImage, error) {

	images, err := r.list(f, p, o)
	if err != nil {
		return nil, err
	}
	return images, nil
}
func (r SQLiteImageRepo) Iterate(f im.FilteringParams, pageSize int) iter.Seq2[im.BaseImage, error] {
	return func(yield func(im.BaseImage, error) bool) {
		var after *im.ImageId
		for {
			page, next, err := r.sliceAfterId(f, pageSize, after)
			if err != nil {
				yield(im.BaseImage{}, err)
				return
			}
			for _, img := range page {
				if !yield(img, nil) {
					return // consumer stopped early
				}
			}
			if len(page) < pageSize {
				return // last page
			}
			after = next
		}
	}
}
func (r SQLiteImageRepo) ImageExistsInCollection(imageId im.ImageId, collectionId clc.CollectionId) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM images_collections WHERE image_id=$1 AND collection_id=$2"
	err := r.Db.QueryRow(query, imageId.String(), collectionId.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking image to collection junction records: %v: %w", err, e.ErrInternal)
	}

	return count > 0, nil
}
func (r SQLiteImageRepo) ImageExists(imageId im.ImageId) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM images WHERE id=$1"
	err := r.Db.QueryRow(query, imageId.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking that image exists: %v: %w", err, e.ErrInternal)
	}

	return count > 0, nil
}
func (r SQLiteImageRepo) GetSpecs(imageId im.ImageId) (*im.ImageSpecs, error) {
	errCtx := "finding image specification"
	var row SpecsRow
	err := r.Db.Get(&row, "SELECT mimetype,width,height,ingested_at FROM images WHERE id = $1", imageId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrNotFound)
		default:
			return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
		}
	}
	return &im.ImageSpecs{MIMEType: row.MIMEType, Width: row.Width, Height: row.Height, IngestedAt: row.IngestedAt}, nil
}
func (r SQLiteImageRepo) AddImage(imageId im.ImageId, hash []byte, specs im.ImageSpecs) error {
	query := "INSERT INTO images (id, hash, mimetype, width, height, ingested_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := r.Db.Exec(query, imageId.String(), hex.EncodeToString(hash), specs.MIMEType,
		specs.Width, specs.Height, specs.IngestedAt)
	if err != nil {
		return fmt.Errorf("inserting image record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteImageRepo) FindImageIdByHash(hash []byte) (*im.ImageId, error) {
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
func (r SQLiteImageRepo) Delete(id im.ImageId) error {
	_, err := r.Db.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleting image record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteImageRepo) RemoveImageFromCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	_, err := r.Db.Exec("DELETE FROM images_collections WHERE image_id = $1 AND collection_id = $2",
		imageId, collectionId)
	if err != nil {
		return fmt.Errorf("removing image from image to collection junction table: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteImageRepo) makeBaseQuery(f im.FilteringParams, pageSize int) sq.SelectBuilder {
	q := sq.StatementBuilder.Select(
		"ic.image_id,ic.collection_id,i.ingested_at,c.name").From(
		"images_collections AS ic").Join(
		"images AS i ON ic.image_id=i.id").Join(
		"collections AS c ON ic.collection_id=c.id")
	q = q.Limit(uint64(pageSize))

	if f.Collection != nil {
		q = q.Where("collection_id=(SELECT id FROM collections WHERE name=?)", *f.Collection)
	}

	return q

}
func (r SQLiteImageRepo) fetchBaseImages(q sq.SelectBuilder) ([]im.BaseImage, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []ListRow{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}
	images := []im.BaseImage{}
	for _, r := range records {
		images = append(images, im.BaseImage{ImageId: r.ImageId, Collection: r.Name})
	}
	return images, nil
}
func (r SQLiteImageRepo) list(f im.FilteringParams, p pa.PaginationParams, o im.OrderingParams) ([]im.BaseImage, error) {
	q := r.makeBaseQuery(f, p.PageSize)
	q = q.Offset((uint64(p.Page-1) * uint64(p.PageSize)))

	if o.IngestTime {
		q = q.OrderBy("i.ingested_at")
	}

	q = q.OrderBy("ic.image_id")
	images, err := r.fetchBaseImages(q)
	if err != nil {
		return nil, err
	}
	return images, nil
}
func (r SQLiteImageRepo) sliceAfterId(f im.FilteringParams, pageSize int, after *im.ImageId) ([]im.BaseImage, *im.ImageId, error) {
	q := r.makeBaseQuery(f, pageSize)
	q = q.OrderBy("ic.image_id")
	if after != nil {
		q = q.Where(sq.Gt{"ic.image_id": after})
	}

	images, err := r.fetchBaseImages(q)
	if err != nil {
		return nil, nil, err
	}
	var next *im.ImageId
	if len(images) > 0 {
		next = &images[len(images)-1].ImageId
	}
	return images, next, nil
}

func NewSQLiteImageRepo(db *sqlx.DB) SQLiteImageRepo {
	return SQLiteImageRepo{Db: db}
}

func NewTestSQLiteImageRepo() SQLiteImageRepo {
	return NewSQLiteImageRepo(s.NewSQLiteDB(":memory:"))
}
