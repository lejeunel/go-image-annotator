package images

import (
	clc "datahub/domain/collections"
	loc "datahub/domain/locations"
	e "datahub/errors"
	g "datahub/generic"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLImageRepo struct {
	Db             *sqlx.DB
	ErrorConverter e.DBErrorConverter
	Logger         *slog.Logger
}

func NewSQLiteImageRepo(db *sqlx.DB, logger *slog.Logger) *SQLImageRepo {

	return &SQLImageRepo{Db: db, ErrorConverter: &e.SQLiteErrorConverter{},
		Logger: logger}

}

type ImageRecord struct {
	Id             ImageId          `db:"id"`
	FileName       string           `db:"filename"`
	CameraId       *loc.CameraId    `db:"camera_id"`
	CapturedAt     time.Time        `db:"captured_at"`
	CreatedAt      time.Time        `db:"created_at"`
	UpdatedAt      time.Time        `db:"updated_at"`
	SHA256         string           `db:"sha256"`
	MIMEType       string           `db:"mimetype"`
	Width          int              `db:"width"`
	Height         int              `db:"height"`
	Type           string           `db:"image_type"`
	CollectionId   clc.CollectionId `db:"collection_id"`
	CollectionName string           `db:"collection_name"`
	Group          string           `db:"group_name"`
}

func FromRecord(record ImageRecord) BaseImage {
	return BaseImage{
		Id:           record.Id,
		FileName:     record.FileName,
		CameraId:     record.CameraId,
		CreatedAt:    record.CreatedAt,
		CapturedAt:   record.CapturedAt,
		UpdatedAt:    record.UpdatedAt,
		SHA256:       record.SHA256,
		MIMEType:     record.MIMEType,
		Width:        record.Width,
		Height:       record.Height,
		Type:         record.Type,
		CollectionId: record.CollectionId,
		Group:        record.Group,
	}
}

func (r *SQLImageRepo) Create(image *Image) (*Image, error) {
	query := `INSERT INTO images
				(id,filename,camera_id,
					captured_at,created_at,updated_at,sha256,
					width,height,mimetype,image_type)
				VALUES
				($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.Db.Exec(query, image.Id, image.FileName,
		image.CameraId,
		image.CapturedAt,
		image.CreatedAt, image.UpdatedAt, image.SHA256, image.Width, image.Height, image.MIMEType,
		image.Type)

	if err != nil {
		return nil, e.ErrDB
	}

	return image, nil
}
func (r *SQLImageRepo) ListWithChecksum(sha256 string) ([]BaseImage, error) {
	query := NewBaseImageQuery()
	query = query.Where("i.sha256=?", sha256)
	query = query.Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, e.ErrDB
	}
	records := []ImageRecord{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		r.Logger.Error("select statement in image listing with checksum", "sql", sql, "args", args, "error", err)
		return nil, e.ErrDB
	}

	if len(records) > 0 {
		images := []BaseImage{}
		for _, record := range records {
			image := FromRecord(record)
			images = append(images, image)
		}
		return images, nil
	}
	return nil, nil

}

func (r *SQLImageRepo) GetBase(imageId ImageId) (*BaseImage, error) {
	record := ImageRecord{}
	query := NewBaseImageQuery().Where("i.id=?", imageId.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, e.ErrDB
	}
	if err := r.Db.Get(&record, sql, args...); err != nil {
		return nil, r.ErrorConverter.Convert(err)
	}
	image := FromRecord(record)

	return &image, nil
}

func (r *SQLImageRepo) Delete(image *Image) error {
	_, err := r.Db.Exec("DELETE FROM images WHERE id=$1", image.Id.String())
	if err != nil {
		return e.ErrDB
	}

	_, err = r.Db.Exec("DELETE FROM annotations WHERE image_id=$1", image.Id.String())
	if err != nil {
		return e.ErrDB
	}
	return nil
}

func (r *SQLImageRepo) Update(image_id ImageId, type_ string, capturedAt time.Time) error {
	_, err := r.Db.Exec("UPDATE images SET image_type=$1, captured_at=$2 WHERE id=$3",
		type_, capturedAt, image_id)

	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLImageRepo) GetAdjacent(currentImage *Image, filters FilterArgs, ordering OrderingArgs, previous bool) (*BaseImage, error) {
	temporalFilter := TemporalFilter{ReferenceImageId: currentImage.Id}
	if ordering.CreatedAt != nil {
		temporalFilter.Field = "created_at"
	} else {
		temporalFilter.Field = "captured_at"
	}
	temporalFilter.Before = previous

	query := NewBaseImageQuery()
	filters.TemporalFilter = &temporalFilter
	query = filters.Apply(query)
	query = query.Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, e.ErrDB
	}

	records := []ImageRecord{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		r.Logger.Error("select statement getting adjacent images", "sql", sql, "args", args, "error", err)
		return nil, e.ErrDB
	}

	if len(records) == 0 {
		return nil, e.ErrNotFound
	}

	images := []BaseImage{}
	for _, record := range records {
		image := FromRecord(record)
		images = append(images, image)
	}

	return &images[0], nil

}

func (r *SQLImageRepo) DeleteImagesInCollection(collection *clc.Collection) error {
	_, err := r.Db.Exec("DELETE FROM image_collection_assoc WHERE collection_id=$1", collection.Id.String())
	if err != nil {
		return e.ErrDB
	}
	return nil
}

func (r *SQLImageRepo) Count(filters FilterArgs) (int64, error) {

	var count int64
	query := NewBaseImageCountQuery()
	query = filters.Apply(query)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, e.ErrDB
	}
	if err := r.Db.QueryRow(sql, args...).Scan(&count); err != nil {
		return 0, e.ErrDB
	}

	return count, nil

}

func (r *SQLImageRepo) List(filters FilterArgs, orderings OrderingArgs, pagination g.PaginationParams) ([]BaseImage, *g.PaginationMeta, error) {

	query := NewBaseImageQuery()
	query = orderings.Apply(query)
	query = filters.Apply(query)
	query = query.Limit(uint64(pagination.Limit())).Offset(uint64(pagination.Offset()))

	records := []ImageRecord{}
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("building SQL statement: %w", e.ErrDB)
	}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		r.Logger.Error("select statement in image listing", "sql", sql, "args", args, "error", err)
		return nil, nil, fmt.Errorf("selecting data: %w", e.ErrDB)
	}

	images := []BaseImage{}
	for _, record := range records {
		image := FromRecord(record)
		images = append(images, image)
	}

	count, err := r.Count(filters)
	if err != nil {
		return nil, nil, err
	}
	paginationMeta := g.NewPaginationMeta(pagination.Page, count, int64(pagination.PageSize))

	return images, &paginationMeta, nil
}

func (r *SQLImageRepo) AssignToCollection(image *Image, collection *clc.Collection) error {
	query := "INSERT INTO image_collection_assoc (id, image_id, collection_id, created_at) VALUES ($1,$2,$3,$4)"
	_, err := r.Db.Exec(query, uuid.New(), image.Id, collection.Id, time.Now())

	if err != nil {
		err = r.ErrorConverter.Convert(err)
		if errors.Is(err, e.ErrDBForeignKeyConstraint) {
			return e.ErrNotFound
		}
	}

	return nil

}

func (r *SQLImageRepo) ImageIsInCollection(image *Image, collection *clc.Collection) (bool, error) {

	var count int64
	query := "SELECT COUNT(*) FROM image_collection_assoc WHERE image_id=$1 AND collection_id=$2"
	if err := r.Db.QueryRow(query, image.Id, collection.Id).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *SQLImageRepo) RemoveImageFromCollection(image *Image) error {
	_, err := r.Db.Exec("DELETE FROM image_collection_assoc WHERE image_id=$1 AND collection_id=$2",
		image.Id, image.Collection.Id)

	if err != nil {
		return err
	}

	return nil

}

func (r *SQLImageRepo) AssignCamera(camera_id loc.CameraId, image_id ImageId) error {
	query := "UPDATE images SET camera_id=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, camera_id, image_id)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLImageRepo) UnassignCamera(id ImageId) error {
	query := "UPDATE images SET camera_id=NULL WHERE id=$1"
	_, err := r.Db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil

}
