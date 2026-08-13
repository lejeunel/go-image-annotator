package scroll

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/modules/scroller"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type SQLiteScrollerRepo struct {
	Db *sqlx.DB
}

type Row struct {
	ImageId      im.ImageId       `db:"image_id"`
	CollectionId clc.CollectionId `db:"collection_id"`
	Collection   string           `db:"name"`
	IngestTime   time.Time        `db:"ingested_at"`
}

func (r SQLiteScrollerRepo) applyScrollOrdering(q sq.SelectBuilder, currentImageId im.ImageId,
	ord im.OrderingArgs, d scroller.ScrollingDirection,
) sq.SelectBuilder {
	for _, orderingCriteria := range ord {
		order := orderingCriteria.Order
		if (order == im.AscOrder) && (d == scroller.ScrollNext) {
			q = q.OrderBy(fmt.Sprintf("i.%v", orderingCriteria.Field))
			q = q.Where(fmt.Sprintf("i.%v>(SELECT %v FROM images WHERE id=?)", orderingCriteria.Field, orderingCriteria.Field), currentImageId)
		}
		if (order == im.AscOrder) && (d == scroller.ScrollPrevious) {
			q = q.OrderBy(fmt.Sprintf("i.%v DESC", orderingCriteria.Field))
			q = q.Where(fmt.Sprintf("i.%v<(SELECT %v FROM images WHERE id=?)", orderingCriteria.Field, orderingCriteria.Field), currentImageId)
		}
		if (order == im.DescOrder) && (d == scroller.ScrollNext) {
			q = q.OrderBy(fmt.Sprintf("i.%v DESC", orderingCriteria.Field))
			q = q.Where(fmt.Sprintf("i.%v<(SELECT %v FROM images WHERE id=?)", orderingCriteria.Field, orderingCriteria.Field), currentImageId)
		}
		if (order == im.DescOrder) && (d == scroller.ScrollPrevious) {
			q = q.OrderBy(fmt.Sprintf("i.%v", orderingCriteria.Field))
			q = q.Where(fmt.Sprintf("i.%v>(SELECT %v FROM images WHERE id=?)", orderingCriteria.Field, orderingCriteria.Field), currentImageId)
		}
	}

	return q
}

func (r SQLiteScrollerRepo) GetAdjacent(id im.ImageId, criteria scroller.ScrollingCriteria,
	d scroller.ScrollingDirection,
) (*im.BaseImage, error) {
	q := sq.StatementBuilder.Select(
		"ic.image_id,ic.collection_id,i.ingested_at,c.name").From(
		"images_collections AS ic").Join(
		"images AS i ON ic.image_id=i.id").Join(
		"collections AS c ON ic.collection_id=c.id")

	if criteria.Collection != nil {
		q = q.Where("c.name=?", *criteria.Collection)
	}

	q = r.applyScrollOrdering(q, id, criteria.OrderingArgs, d)
	q = q.Limit(1)
	sqlQuery, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	var row Row
	if err := r.Db.Get(&row, sqlQuery, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	result := im.BaseImage{ImageId: row.ImageId, Collection: row.Collection}
	return &result, nil
}

func (r SQLiteScrollerRepo) ImageMustExist(id im.ImageId) error {
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

func (r SQLiteScrollerRepo) CollectionMustExist(collection string) error {
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

func NewSQLiteScrollerRepo(db *sqlx.DB) SQLiteScrollerRepo {
	return SQLiteScrollerRepo{Db: db}
}
