package annotation

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	c "github.com/lejeunel/go-image-annotator/entities/collection"
	i "github.com/lejeunel/go-image-annotator/entities/image"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	sl "github.com/lejeunel/go-image-annotator/infra/db/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type SQLiteAnnotationRepo struct {
	Db *sqlx.DB
	sl.SQLiteLabelRepo
}

type AnnotationRow struct {
	Id          a.AnnotationId `db:"id"`
	LabelId     l.LabelId      `db:"label_id"`
	Type        string         `db:"type"`
	Coordinates string         `db:"coordinates"`
}

type BoundingBoxSpecs struct {
	Xc     float32 `json:"xc"`
	Yc     float32 `json:"yc"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

func (r *SQLiteAnnotationRepo) AddImageLabel(imageId i.ImageId, collectionId c.CollectionId, ann a.ImageLabel) error {
	query := "INSERT INTO annotations (id, image_id, collection_id, label_id, type) VALUES ($1,$2,$3,$4,$5)"
	_, err := r.Db.Exec(query, ann.Id, imageId, collectionId, ann.Label.Id, "image")
	if err != nil {
		return fmt.Errorf("adding image label annotation record: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r *SQLiteAnnotationRepo) findLabelById(labelId l.LabelId) (*l.Label, error) {

	rec := sl.LabelRecord{}
	err := r.Db.Get(&rec,
		"SELECT id,name,description FROM labels WHERE id=$1", labelId)
	if err != nil {
		return nil, fmt.Errorf("fetching label by id %v: %w", labelId, e.ErrInternal)
	}
	return &l.Label{Id: rec.Id, Name: rec.Name, Description: rec.Description}, nil

}
func (r *SQLiteAnnotationRepo) FindImageLabels(imageId i.ImageId, collectionId c.CollectionId) ([]a.ImageLabel, error) {
	query := "SELECT id,label_id,type FROM annotations WHERE image_id=$1 AND collection_id=$2 AND type='image'"

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
		imageLabels = append(imageLabels, a.ImageLabel{Id: rec.Id, Label: *label})
	}

	return imageLabels, nil
}
func (r *SQLiteAnnotationRepo) RemoveAnnotation(id a.AnnotationId) error {
	_, err := r.Db.Exec("DELETE FROM annotations WHERE id=$1", id)

	if err != nil {
		return fmt.Errorf("deleting annotation record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteAnnotationRepo) RemoveImageLabel(imageId i.ImageId, collectionId c.CollectionId, labelId l.LabelId) error {
	_, err := r.Db.Exec("DELETE FROM annotations WHERE image_id=$1 AND collection_id=$2 AND label_id=$3 AND type='image'",
		imageId, collectionId, labelId)

	if err != nil {
		return fmt.Errorf("deleting image label: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteAnnotationRepo) AddBoundingBox(imageId i.ImageId, collectionId c.CollectionId, box a.BoundingBox) error {

	coordsBytes, _ := json.Marshal(BoundingBoxSpecs{Xc: box.Xc, Yc: box.Yc, Width: box.Width, Height: box.Height})
	coordsString := string(coordsBytes)
	query := "INSERT INTO annotations (id, image_id, collection_id, label_id, type, coordinates) VALUES ($1,$2,$3,$4,$5,$6)"
	_, err := r.Db.Exec(query, box.Id, imageId, collectionId, box.Label.Id, "bounding_box", coordsString)
	if err != nil {
		return fmt.Errorf("inserting record in annotation table: %v: %w", err, e.ErrInternal)
	}

	return nil
}
func (r *SQLiteAnnotationRepo) FindBoundingBoxes(imageId i.ImageId, collectionId c.CollectionId) ([]a.BoundingBox, error) {
	query := "SELECT id,label_id,type,coordinates FROM annotations WHERE image_id=$1 AND collection_id=$2 AND type='bounding_box'"

	errCtx := "querying image annotations"
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
		boxes = append(boxes,
			a.NewBoundingBox(rec.Id, specs.Xc, specs.Yc, specs.Width, specs.Height, *label))
	}

	return boxes, nil
}
func (r *SQLiteAnnotationRepo) UpdateLabelOfAnnotation(id a.AnnotationId, labelId l.LabelId) error {
	query := "UPDATE annotations SET label_id=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, labelId, id)

	if err != nil {
		return fmt.Errorf("updating bounding box label: %v: %w", err, e.ErrInternal)
	}

	return nil

}
func (r *SQLiteAnnotationRepo) UpdateBoundingBoxCoordinates(id a.AnnotationId, xc, yc, width, height float32) error {
	errCtx := "updating bounding box coordinates"
	if err := a.ValidateBoundingBox(xc, yc, width, height); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}

	coordsBytes, _ := json.Marshal(BoundingBoxSpecs{Xc: xc, Yc: yc, Width: width, Height: height})
	coordsString := string(coordsBytes)
	query := "UPDATE annotations SET coordinates=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, coordsString, id)
	if err != nil {
		return fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteAnnotationRepo) UpdateBoundingBox(id a.AnnotationId, u a.BoundingBoxUpdatables) error {
	if err := r.UpdateLabelOfAnnotation(id, u.LabelId); err != nil {
		return err
	}

	if err := r.UpdateBoundingBoxCoordinates(id, u.Xc, u.Yc, u.Width, u.Height); err != nil {
		return err
	}
	return nil
}
func (r *SQLiteAnnotationRepo) GroupOfAnnotation(id a.AnnotationId) (*string, error) {
	var group string
	err := r.Db.Get(&group,
		`SELECT "group" FROM collections WHERE id=(SELECT collection_id FROM annotations WHERE id=$1 )`, id)
	if err != nil {
		return nil, fmt.Errorf("fetching group of annotation by id %v: %w", id, e.ErrInternal)
	}
	return &group, nil

}
func NewSQLiteAnnotationRepo(db *sqlx.DB) *SQLiteAnnotationRepo {
	return &SQLiteAnnotationRepo{Db: db, SQLiteLabelRepo: *sl.NewSQLiteLabelRepo(db)}
}
