package scroll

import (
	"fmt"

	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lejeunel/go-image-annotator-v2/app/annotator/scroller"
	clc "github.com/lejeunel/go-image-annotator-v2/entities/collection"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

type SQLiteScrollerRepo struct {
	Db *sqlx.DB
}

type Row struct {
	ImageId      im.ImageId       `db:"image_id"`
	CollectionId clc.CollectionId `db:"collection_id"`
}

func (r *SQLiteScrollerRepo) GetAdjacent(id im.ImageId, criteria scroller.ScrollingCriteria,
	d scroller.ScrollingDirection) (*im.BaseImage, error) {
	q := sq.StatementBuilder.Select("id").From("images")

	result := im.BaseImage{}

	if d == scroller.ScrollNext {
		q = q.Where(fmt.Sprintf("id>'%v'", id))
		q = q.OrderBy("id")
	} else {
		q = q.Where(fmt.Sprintf("id<'%v'", id))
		q = q.OrderBy("id DESC")
	}

	if criteria.Collection != nil {
		q = q.Where(fmt.Sprintf("id IN (SELECT image_id FROM images_collections WHERE collection_id=(SELECT id FROM collections WHERE name='%v'))",
			*criteria.Collection))
		result.Collection = *criteria.Collection
	}

	q = q.Limit(1)

	sqlQuery, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	var adjId im.ImageId
	if err := r.Db.Get(&adjId, sqlQuery, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	result.ImageId = adjId
	return &result, nil
}

func (r *SQLiteScrollerRepo) ImageMustExist(id im.ImageId) error {
	var count int64
	query := "SELECT COUNT(*) FROM images WHERE id=$1"
	err := r.Db.QueryRow(query, id.String()).Scan(&count)
	if err != nil {
		return fmt.Errorf("checking whether image exists: %w: %w", err, e.ErrInternal)
	}
	if count == 0 {
		return fmt.Errorf("checking whether image exists: %w", e.ErrNotFound)
	}
	return nil
}

func (r *SQLiteScrollerRepo) CollectionMustExist(collection string) error {
	var count int64
	query := "SELECT COUNT(*) FROM collections WHERE name=$1"
	err := r.Db.QueryRow(query, collection).Scan(&count)
	if err != nil {
		return fmt.Errorf("checking whether collection exists: %w: %w", err, e.ErrInternal)
	}
	if count == 0 {
		return fmt.Errorf("checking whether collection exists: %w", e.ErrNotFound)
	}
	return nil
}

func NewSQLiteScrollerRepo(db *sqlx.DB) *SQLiteScrollerRepo {
	return &SQLiteScrollerRepo{Db: db}
}
