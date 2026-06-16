package annotation

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	sl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	c "github.com/lejeunel/go-image-annotator/entities/collection"
	i "github.com/lejeunel/go-image-annotator/entities/image"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type SQLiteAnnotationRepo struct {
	Db *sqlx.DB
}

type AnnotationRow struct {
	Id          a.AnnotationId `db:"id"`
	LabelId     l.LabelId      `db:"label_id"`
	Type        string         `db:"type"`
	Coordinates string         `db:"coordinates"`
	Author      *u.UserId      `db:"author"`
	Time        *time.Time     `db:"touched_at"`
}

type BoundingBoxSpecs struct {
	Xc     float32 `json:"xc"`
	Yc     float32 `json:"yc"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
	Angle  float32 `json:"angle"`
}

func (r SQLiteAnnotationRepo) AddImageLabel(imageId i.ImageId, collectionId c.CollectionId, ann a.ImageLabel, userId *u.UserId, t *time.Time) error {
	query := "INSERT INTO annotations (id, image_id, collection_id, label_id, type, author, touched_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	_, err := r.Db.Exec(query, ann.Id, imageId, collectionId, ann.Label.Id, "image", userId, t)
	if err != nil {
		return fmt.Errorf("adding image label annotation record: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteAnnotationRepo) findLabelById(labelId l.LabelId) (*l.Label, error) {

	rec := sl.LabelRecord{}
	err := r.Db.Get(&rec,
		"SELECT id,name,description FROM labels WHERE id=$1", labelId)
	if err != nil {
		return nil, fmt.Errorf("fetching label by id %v: %w", labelId, e.ErrInternal)
	}
	return &l.Label{Id: rec.Id, Name: rec.Name, Description: rec.Description}, nil

}
func (r SQLiteAnnotationRepo) FindImageLabels(imageId i.ImageId, collectionId c.CollectionId) ([]a.ImageLabel, error) {
	query := "SELECT id,label_id,type,author,touched_at FROM annotations WHERE image_id=$1 AND collection_id=$2 AND type='image'"

	errCtx := "querying image annotations"
	records := []AnnotationRow{}
	if err := r.Db.Select(&records, query, imageId, collectionId); err != nil {
		return nil, fmt.Errorf("%v: applying query: %v: %w", errCtx, err, e.ErrInternal)
	}

	imageLabels := []a.ImageLabel{}
	for _, rec := range records {
		label, err := r.findLabelById(rec.LabelId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", errCtx, err)
		}
		imageLabels = append(imageLabels, a.ImageLabel{Id: rec.Id, Label: *label, Author: rec.Author, Time: rec.Time})
	}

	return imageLabels, nil
}
func (r SQLiteAnnotationRepo) RemoveAnnotation(id a.AnnotationId) error {
	_, err := r.Db.Exec("DELETE FROM annotations WHERE id=$1", id)

	if err != nil {
		return fmt.Errorf("deleting annotation record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteAnnotationRepo) RemoveImageLabel(imageId i.ImageId, collectionId c.CollectionId, labelId l.LabelId) error {
	_, err := r.Db.Exec("DELETE FROM annotations WHERE image_id=$1 AND collection_id=$2 AND label_id=$3 AND type='image'",
		imageId, collectionId, labelId)

	if err != nil {
		return fmt.Errorf("deleting image label: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteAnnotationRepo) AddBoundingBox(imageId i.ImageId, collectionId c.CollectionId, box a.BoundingBox, userId *u.UserId, t *time.Time) error {

	coordsBytes, _ := json.Marshal(BoundingBoxSpecs{Xc: box.Xc, Yc: box.Yc, Width: box.Width, Height: box.Height, Angle: box.Angle})
	coordsString := string(coordsBytes)
	query := "INSERT INTO annotations (id, image_id, collection_id, label_id, type, coordinates, author, touched_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"
	_, err := r.Db.Exec(query, box.Id, imageId, collectionId, box.Label.Id, "bounding_box", coordsString, userId, t)
	if err != nil {
		return fmt.Errorf("inserting bounding box: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteAnnotationRepo) FindBoundingBoxes(imageId i.ImageId, collectionId c.CollectionId) ([]a.BoundingBox, error) {
	query := "SELECT id,label_id,type,coordinates,author,touched_at FROM annotations WHERE image_id=$1 AND collection_id=$2 AND type='bounding_box'"

	errCtx := "querying bounding-box annotations"
	records := []AnnotationRow{}
	if err := r.Db.Select(&records, query, imageId, collectionId); err != nil {
		return nil, fmt.Errorf("%v: applying query: %v: %w", errCtx, err, e.ErrInternal)
	}

	boxes := []a.BoundingBox{}
	for _, rec := range records {
		var specs BoundingBoxSpecs
		err := json.Unmarshal([]byte(rec.Coordinates), &specs)
		if err != nil {
			return nil, fmt.Errorf("%v: unmarshaling bounding box specs: %+v: %w: %w",
				errCtx, rec.Coordinates, err, e.ErrInternal)
		}
		label, err := r.findLabelById(rec.LabelId)
		box := a.NewBoundingBox(rec.Id, specs.Xc, specs.Yc, specs.Width, specs.Height, *label,
			a.WithAngle(specs.Angle))
		if rec.Author != nil {
			box.Author = rec.Author
		}
		if rec.Time != nil {
			box.Time = rec.Time
		}
		boxes = append(boxes, box)
	}

	return boxes, nil
}
func (r SQLiteAnnotationRepo) UpdateLabelOfAnnotation(id a.AnnotationId, labelId l.LabelId, userId *u.UserId, t *time.Time) error {
	errCtx := "updating bounding box"
	if err := r.UpdateAuthor(id, userId); err != nil {
		return fmt.Errorf("%v: updating author: %w", errCtx, err)
	}
	if err := r.UpdateTime(id, t); err != nil {
		return fmt.Errorf("%v: updating time: %w", errCtx, err)
	}

	query := "UPDATE annotations SET label_id=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, labelId, id)

	if err != nil {
		return fmt.Errorf("updating bounding box label: %v: %w", err, e.ErrInternal)
	}

	return nil

}
func (r SQLiteAnnotationRepo) UpdateBoundingBoxCoordinates(id a.AnnotationId, xc, yc, width, height, angle float32) error {
	errCtx := "updating bounding box coordinates"
	if err := a.ValidateBoundingBox(xc, yc, width, height, angle); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}

	coordsBytes, _ := json.Marshal(BoundingBoxSpecs{Xc: xc, Yc: yc, Width: width, Height: height, Angle: angle})
	coordsString := string(coordsBytes)
	query := "UPDATE annotations SET coordinates=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, coordsString, id)
	if err != nil {
		return fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteAnnotationRepo) UpdateAuthor(id a.AnnotationId, userId *u.UserId) error {
	query := "UPDATE annotations SET author=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, userId, id)

	if err != nil {
		return fmt.Errorf("%w: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteAnnotationRepo) UpdateTime(id a.AnnotationId, t *time.Time) error {
	query := "UPDATE annotations SET touched_at=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, t, id)

	if err != nil {
		return fmt.Errorf("%w: %w", err, e.ErrInternal)
	}

	return nil
}
func (r SQLiteAnnotationRepo) UpdateBoundingBox(id a.AnnotationId, u a.BoundingBoxUpdatables, userId *u.UserId, t *time.Time) error {
	errCtx := "updating bounding box"
	if err := r.UpdateLabelOfAnnotation(id, u.LabelId, userId, t); err != nil {
		return fmt.Errorf("%v: updating label: %w", errCtx, err)
	}

	if err := r.UpdateBoundingBoxCoordinates(id, u.Xc, u.Yc, u.Width, u.Height, u.Angle); err != nil {
		return fmt.Errorf("%v: updating coordinates: %w", errCtx, err)
	}
	return nil
}
func (r SQLiteAnnotationRepo) GroupOfAnnotation(id a.AnnotationId) (*string, error) {
	var group string
	err := r.Db.Get(&group,
		`SELECT name FROM groups WHERE id=(SELECT group_id FROM collections WHERE id=(SELECT collection_id FROM annotations WHERE id=$1))`,
		id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetching group of annotation by id %v: %w", id, e.ErrInternal)
	}
	return &group, nil
}
func NewSQLiteAnnotationRepo(db *sqlx.DB) SQLiteAnnotationRepo {
	return SQLiteAnnotationRepo{Db: db}
}
