package image

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"go.tomakado.io/dumbql/schema"
	"iter"
	"time"

	sq "github.com/Masterminds/squirrel"
	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	qu "github.com/lejeunel/go-image-annotator/modules/query"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	ss "github.com/lejeunel/go-image-annotator/shared/sql"
)

func MakeQueryParsers() (qu.FilterParser, qu.OrderingStrConverter) {
	b := schema.NewSchemaBuilder()
	b.AddField("collection", schema.Is[string]())
	b.AddField("ingested_at", schema.Is[string]())
	filteringStrConverter := qu.NewFilterParser(b.Build(), qu.WithFieldNameMapping("collection", "collections.name"))
	orderingStrConverter := qu.NewOrderingConverter(qu.WithOrderingField("collection"), qu.WithOrderingField("ingested_at"))
	return filteringStrConverter, orderingStrConverter
}

type FilterStrParser interface {
	Parse(string) (ss.SQLizer, error)
}

type OrderStrParser interface {
	Parse(string) (string, error)
}

type SQLiteImageRepo struct {
	Db adb.Querier
	FilterStrParser
	OrderStrParser
}

func NewSQLiteImageRepo(db adb.Querier, fp FilterStrParser, op OrderStrParser) SQLiteImageRepo {
	return SQLiteImageRepo{db, fp, op}
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

func (r SQLiteImageRepo) AddToCollection(imageId im.ImageId, collection clc.CollectionName) error {
	query := "INSERT INTO images_collections (image_id, collection_id) VALUES ($1,(SELECT id FROM collections WHERE name=$2))"
	_, err := r.Db.Exec(query, imageId.String(), collection)
	if err != nil {
		return fmt.Errorf("inserting image record into junction table: %v: %w", err, e.ErrInternal)
	}

	return nil
}

func (r SQLiteImageRepo) Count(f im.FilterQueryStr) (*int64, error) {
	expr, err := r.FilterStrParser.Parse(f)
	if err != nil {
		return nil, fmt.Errorf("parsing filtering string %v: %v: %w", f, err, e.ErrInternal)
	}
	q := sq.StatementBuilder.Select("COUNT(*)")
	q = q.From("images_collections").Join("collections ON collections.id = images_collections.collection_id")
	sql, args, err := q.Where(expr).ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query from filters %v: %v: %w", f, err, e.ErrInternal)
	}

	var count int64
	if err := r.Db.Get(&count, sql, args...); err != nil {
		return nil, fmt.Errorf("querying using filters %v: %v: %w", f, err, e.ErrInternal)
	}

	return &count, nil
}

func (r SQLiteImageRepo) Slice(
	f im.FilterQueryStr,
	p pa.PaginationParams,
	o im.OrderingStr,
) ([]im.BaseImage, error) {
	q, err := r.makeBaseQuery(f, p.PageSize)
	if err != nil {
		return nil, err
	}
	qq := *q
	qq = qq.Offset((uint64(p.Page-1) * uint64(p.PageSize)))

	if o != "" {
		orderStr, err := r.OrderStrParser.Parse(o)
		if err != nil {
			return nil, err
		}
		qq = qq.OrderBy(orderStr)

	}
	qq = qq.OrderBy("ic.image_id")
	images, err := r.fetchBaseImages(qq)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (r SQLiteImageRepo) sliceAfterId(
	f im.FilterQueryStr,
	pageSize int,
	after *im.ImageId,
) ([]im.BaseImage, *im.ImageId, error) {
	q0, err := r.makeBaseQuery(f, pageSize)
	if err != nil {
		return nil, nil, err
	}
	q1 := q0.OrderBy("ic.image_id")
	if after != nil {
		q1 = q1.Where(sq.Gt{"ic.image_id": after})
	}

	images, err := r.fetchBaseImages(q1)
	if err != nil {
		return nil, nil, err
	}
	var next *im.ImageId
	if len(images) > 0 {
		next = &images[len(images)-1].ImageId
	}
	return images, next, nil
}

func (r SQLiteImageRepo) Iterate(f im.FilterQueryStr, pageSize int) iter.Seq2[im.BaseImage, error] {
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

func (r SQLiteImageRepo) ImageExistsInCollection(
	imageId im.ImageId,
	collection clc.CollectionName,
) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM images_collections WHERE image_id=$1 AND collection_id=(SELECT id FROM collections WHERE name=$2)"
	err := r.Db.QueryRow(query, imageId.String(), collection).Scan(&count)
	if err != nil {
		return false, fmt.Errorf(
			"checking image to collection junction records: %v: %w",
			err,
			e.ErrInternal,
		)
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

func (r SQLiteImageRepo) GetSpecs(imageId im.ImageId) (*im.Specs, error) {
	errCtx := "finding image specification"
	var row SpecsRow
	err := r.Db.Get(
		&row,
		"SELECT mimetype,width,height,ingested_at FROM images WHERE id = $1",
		imageId,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrNotFound)
		default:
			return nil, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
		}
	}
	return &im.Specs{
		MIMEType:   row.MIMEType,
		Width:      row.Width,
		Height:     row.Height,
		IngestedAt: row.IngestedAt,
	}, nil
}

func (r SQLiteImageRepo) AddImage(imageId im.ImageId, hash []byte, specs im.Specs) error {
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

func (r SQLiteImageRepo) RemoveImageFromCollection(
	imageId im.ImageId,
	collection clc.CollectionName,
) error {
	_, err := r.Db.Exec("DELETE FROM images_collections WHERE image_id = $1 AND collection_id = (SELECT id FROM collections WHERE name=$2)",
		imageId, collection)
	if err != nil {
		return fmt.Errorf(
			"removing image from image to collection junction table: %v: %w",
			err,
			e.ErrInternal,
		)
	}
	return nil
}

func (r SQLiteImageRepo) makeBaseQuery(f im.FilterQueryStr, pageSize int) (*sq.SelectBuilder, error) {
	q := sq.StatementBuilder.Select(
		"ic.image_id,ic.collection_id,i.ingested_at,c.name").From(
		"images_collections AS ic").Join(
		"images AS i ON ic.image_id=i.id").Join(
		"collections AS c ON ic.collection_id=c.id")
	q = q.Limit(uint64(pageSize))

	if f != "" {
		expr, err := r.FilterStrParser.Parse(f)
		if err != nil {
			return nil, err
		}
		q.Where(expr)

	}

	return &q, nil
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

func (r SQLiteImageRepo) IsUsed(id im.ImageId) (*bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM images_collections WHERE image_id=$1"
	err := r.Db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf(
			"counting number of collections using image %v: %v: %w",
			id,
			err,
			e.ErrInternal,
		)
	}
	var isUsed bool
	if count > 0 {
		isUsed = true
	}
	return &isUsed, nil
}
