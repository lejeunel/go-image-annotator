package image

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"iter"
	"strings"
	"time"

	"go.tomakado.io/dumbql/query"
	"go.tomakado.io/dumbql/schema"

	sq "github.com/Masterminds/squirrel"
	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	qu "github.com/lejeunel/go-image-annotator/modules/query"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

func makeOrderingWindowExpr(parser OrderStrParser, function string, ordering im.OrderStr, outName string) (*string, error) {
	if ordering == "" {
		ordering = "image_id"
	}
	args, err := parser.Parse(ordering)
	if err != nil {
		return nil, err
	}

	var terms []string
	for _, a := range args {
		if a.Order == im.DescOrder {
			terms = append(terms, a.Field+" "+"DESC")
		} else {
			terms = append(terms, a.Field)
		}
	}

	res := function + " OVER (ORDER BY " + strings.Join(terms, ",") + ") " + outName
	return &res, nil
}

func MakeQueryParsers() (qu.FilterParser, qu.OrderParser) {
	sb := schema.NewSchemaBuilder()
	sb.AddField("collection", schema.Is[string]())
	sb.AddField("ingested_at", schema.Is[string]())
	sb.AddRegExpField(`^meta\..*$`, schema.Any(schema.Is[float64](), schema.Is[string](), schema.Is[bool]()))

	rb := query.NewRenamerBuilder()
	// rb.Add(`\bcollection\b`, `collections.name`)
	// rb.Add(`\bingested_at\b`, `i.ingested_at`)
	rb.Add(`\bmeta\.(.*)\b`, `json_extract(m.meta, '$.$1')`)

	filterParser := qu.NewFilterParser(
		sb.Build(),
		qu.WithRenamer(rb.Build()),
	)
	ob := qu.NewOrderParserBuilder()
	ob.AddField("image_id")
	ob.AddField("ingested_at")
	ob.AddRegExpField(`^meta\..*$`)
	ob.AddRenameRule(`\bmeta\.(.*)\b`, `json_extract(m.meta, '$.$1')`)
	// rb.Add(`\bingested_at\b`, `i.ingested_at`)
	// ob.AddRenameRule(`\bimage_id\b`, `i.id`)
	// ob.AddRenameRule(`\bingested_at\b`, `i.ingested_at`)
	// ob.AddRenameRule(`\bcollection\b`, `collections.name`)

	orderParser := ob.Build()

	return filterParser, orderParser
}

type OrderStrParser interface {
	Parse(string) (im.OrderingArgs, error)
}

type FilterParser interface {
	ParseToSql(q string) (*qu.SQLizer, error)
}

type ImageRepo struct {
	Db adb.Querier
	FilterParser
	OrderStrParser
}

func NewImageRepo(db adb.Querier, fp FilterParser, op OrderStrParser) ImageRepo {
	return ImageRepo{db, fp, op}
}

type Row struct {
	ImageId        im.ImageId         `db:"id"`
	CollectionName clc.CollectionName `db:"name"`
}

type AdjacencyRow struct {
	ImageId        *im.ImageId         `db:"image_id"`
	Collection     string              `db:"collection"`
	IngestedAt     time.Time           `db:"ingested_at"`
	PrevId         *im.ImageId         `db:"prev_id"`
	PrevCollection *clc.CollectionName `db:"prev_collection"`
	NextId         *im.ImageId         `db:"next_id"`
	NextCollection *clc.CollectionName `db:"next_collection"`
}

type SpecsRow struct {
	MIMEType   string    `db:"mimetype"`
	Width      int       `db:"width"`
	Height     int       `db:"height"`
	IngestedAt time.Time `db:"ingested_at"`
}

func (r ImageRepo) AddToCollection(imageId im.ImageId, collection clc.CollectionName) error {
	query := "INSERT INTO images_collections (image_id, collection_id) VALUES ($1,(SELECT id FROM collections WHERE name=$2))"
	_, err := r.Db.Exec(query, imageId.String(), collection)
	if err != nil {
		return fmt.Errorf("inserting image record into junction table: %v: %w", err, e.ErrInternal)
	}

	return nil
}

func (r ImageRepo) Count(f im.FilterStr) (*int64, error) {
	errCtx := fmt.Errorf("counting using filters %v", f)
	q := sq.StatementBuilder.Select("COUNT(*)")
	q = q.From("images_collections AS ic").Join(
		"images AS i ON ic.image_id=i.id").Join(
		"collections ON ic.collection_id=collections.id").LeftJoin(
		"metadata AS m ON ic.image_id=m.image_id AND ic.collection_id=m.collection_id")

	if f != "" {
		sqlizer, err := r.FilterParser.ParseToSql(f)
		if err != nil {
			return nil, err
		}
		q = q.Where(*sqlizer)
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, e.ErrInternal)
	}

	var count int64
	if err := r.Db.Get(&count, sql, args...); err != nil {
		return nil, fmt.Errorf("querying using filters %v: %v: %w", f, err, e.ErrInternal)
	}

	return &count, nil
}
func (r ImageRepo) applyOrderingStr(q sq.SelectBuilder, o im.OrderStr) (*sq.SelectBuilder, error) {
	if o != "" {
		args, err := r.OrderStrParser.Parse(o)
		if err != nil {
			return nil, err
		}
		for _, a := range args {
			if a.Order == im.DescOrder {
				q = q.OrderBy(a.Field + " " + "DESC")
			} else {
				q = q.OrderBy(a.Field)
			}
		}
	}
	return &q, nil
}

func (r ImageRepo) Slice(
	f im.FilterStr,
	p pa.PaginationParams,
	o im.OrderStr,
) ([]im.BaseImage, error) {
	q := r.makeBaseSelectQuery()
	qf, err := r.applyFilters(q, f)
	if err != nil {
		return nil, err
	}

	qfo, err := r.applyOrderingStr(*qf, o)
	if err != nil {
		return nil, err
	}

	qs := *qfo
	qs = qs.Limit(uint64(p.PageSize))
	qs = qs.Offset((uint64(p.Page-1) * uint64(p.PageSize)))
	qs = qs.OrderBy("ic.image_id")
	images, err := r.fetchBaseImages(qs)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (r ImageRepo) sliceAfterId(
	f im.FilterStr,
	pageSize int,
	after *im.ImageId,
) ([]im.BaseImage, *im.ImageId, error) {
	q := r.makeBaseSelectQuery()
	qf, err := r.applyFilters(q, f)
	if err != nil {
		return nil, nil, err
	}

	qs := *qf
	qs = qs.OrderBy("ic.image_id")
	qs = qs.Limit(uint64(pageSize))
	if after != nil {
		qs = qs.Where(sq.Gt{"ic.image_id": after})
	}

	images, err := r.fetchBaseImages(qs)
	if err != nil {
		return nil, nil, err
	}
	var next *im.ImageId
	if len(images) > 0 {
		next = &images[len(images)-1].ImageId
	}
	return images, next, nil
}

func (r ImageRepo) Iterate(f im.FilterStr, pageSize int) iter.Seq2[im.BaseImage, error] {
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

func (r ImageRepo) ImageExistsInCollection(
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

func (r ImageRepo) ImageExists(imageId im.ImageId) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM images WHERE id=$1"
	err := r.Db.QueryRow(query, imageId.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking that image exists: %v: %w", err, e.ErrInternal)
	}

	return count > 0, nil
}

func (r ImageRepo) GetSpecs(imageId im.ImageId) (*im.Specs, error) {
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

func (r ImageRepo) AddImage(imageId im.ImageId, hash []byte, specs im.Specs) error {
	query := "INSERT INTO images (id, hash, mimetype, width, height, ingested_at) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := r.Db.Exec(query, imageId.String(), hex.EncodeToString(hash), specs.MIMEType,
		specs.Width, specs.Height, specs.IngestedAt)
	if err != nil {
		return fmt.Errorf("inserting image record: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r ImageRepo) FindImageIdByHash(hash []byte) (*im.ImageId, error) {
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

func (r ImageRepo) Delete(id im.ImageId) error {
	_, err := r.Db.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleting image record: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r ImageRepo) RemoveImageFromCollection(
	imageId im.ImageId,
	collection clc.CollectionName,
) error {
	_, err := r.Db.Exec(
		"DELETE FROM images_collections WHERE image_id = $1 AND collection_id = (SELECT id FROM collections WHERE name=$2)",
		imageId,
		collection,
	)
	if err != nil {
		return fmt.Errorf(
			"removing image from image to collection junction table: %v: %w",
			err,
			e.ErrInternal,
		)
	}
	return nil
}

func (r ImageRepo) applyFilters(q sq.SelectBuilder, f im.FilterStr) (*sq.SelectBuilder, error) {
	if f != "" {
		sqlizer, err := r.FilterParser.ParseToSql(f)
		if err != nil {
			return nil, err
		}
		q = q.Where(*sqlizer)
		return &q, nil
	}

	return &q, nil

}

func (r ImageRepo) makeBaseSelectQuery() sq.SelectBuilder {
	return sq.StatementBuilder.Select(
		"i.id,collections.name").From(
		"images_collections AS ic").Join(
		"images AS i ON ic.image_id=i.id").Join(
		"collections ON ic.collection_id=collections.id").LeftJoin(
		"metadata AS m ON ic.image_id=m.image_id AND ic.collection_id=m.collection_id")

}

func (r ImageRepo) fetchBaseImages(q sq.SelectBuilder) ([]im.BaseImage, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Row{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}
	images := []im.BaseImage{}
	for _, r := range records {
		images = append(images, im.BaseImage{ImageId: r.ImageId, Collection: r.CollectionName})
	}
	return images, nil
}

func (r ImageRepo) IsUsed(id im.ImageId) (*bool, error) {
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

func (r ImageRepo) GetAdjacent(
	currentId im.ImageId,
	currentCollection string,
	f im.FilterStr,
	o im.OrderStr,
	d im.ScrollingDirection,
) (*im.BaseImage, error) {

	prevIdWinExpr, innerCountErr := makeOrderingWindowExpr(r.OrderStrParser, "LAG(image_id, 1)", o, "prev_id")
	if innerCountErr != nil {
		return nil, innerCountErr
	}
	prevCollectionWinExpr, innerCountErr := makeOrderingWindowExpr(r.OrderStrParser, "LAG(collection, 1)", o, "prev_collection")
	if innerCountErr != nil {
		return nil, innerCountErr
	}
	nextIdWinExpr, innerCountErr := makeOrderingWindowExpr(r.OrderStrParser, "LEAD(image_id, 1)", o, "next_id")
	if innerCountErr != nil {
		return nil, innerCountErr
	}
	nextCollectionWinExpr, innerCountErr := makeOrderingWindowExpr(r.OrderStrParser, "LEAD(collection, 1)", o, "next_collection")
	if innerCountErr != nil {
		return nil, innerCountErr
	}
	all := sq.StatementBuilder.Select(
		"i.id AS image_id",
		"i.ingested_at AS ingested_at",
		"collections.name AS collection",
	).From(
		"images_collections AS ic").Join(
		"images AS i ON ic.image_id=i.id").Join(
		"collections ON ic.collection_id=collections.id")
	filtered, err := r.applyFilters(all, f)
	if err != nil {
		return nil, err
	}
	ordered, err := r.applyOrderingStr(*filtered, o)
	if err != nil {
		return nil, err
	}

	adj := sq.StatementBuilder.Select(
		"image_id AS image_id",
		"collection AS collection",
		"ingested_at",
		*prevIdWinExpr,
		*prevCollectionWinExpr,
		*nextIdWinExpr,
		*nextCollectionWinExpr,
	).FromSelect(*ordered, "adjacency")

	adjSelf := sq.StatementBuilder.Select(
		"image_id AS image_id",
		"collection AS collection",
		"ingested_at",
		"prev_id",
		"prev_collection",
		"next_id",
		"next_collection",
	).FromSelect(adj, "adj").Where(
		"image_id = ? AND collection= ?", currentId, currentCollection,
	)

	adjSql, adjArgs, err := adjSelf.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	var row AdjacencyRow

	err = r.Db.Get(&row, adjSql, adjArgs...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	var res im.BaseImage
	fmt.Printf("got adjacencies %+v\n", row)
	if d == im.ScrollNext {
		if row.NextId == nil {
			return nil, nil
		}
		res.ImageId = *row.NextId
		res.Collection = *row.NextCollection

	} else {
		if row.PrevId == nil {
			return nil, nil
		}
		res.ImageId = *row.PrevId
		res.Collection = *row.PrevCollection

	}

	return &res, nil
}
