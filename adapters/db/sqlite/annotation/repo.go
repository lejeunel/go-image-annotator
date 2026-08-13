package annotation

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	adb "github.com/lejeunel/go-image-annotator/adapters/db"
	sl "github.com/lejeunel/go-image-annotator/adapters/db/sqlite/label"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	c "github.com/lejeunel/go-image-annotator/entities/collection"
	i "github.com/lejeunel/go-image-annotator/entities/image"
	l "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type AnnotationRepo struct {
	Db adb.Querier
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

type PointSpec struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type PolygonSpecs struct {
	Points []PointSpec `json:"points"`
}

func (r AnnotationRepo) AddImageLabel(
	imageId i.ImageId,
	collection c.CollectionName,
	ann a.ImageLabel,
	userId *u.UserId,
	t *time.Time,
) error {
	query := `INSERT INTO annotations (id, image_id, collection_id, label_id, type, author, touched_at)
		VALUES ($1,$2,(SELECT id FROM collections WHERE name=$3),$4,$5,$6,$7)`
	_, err := r.Db.Exec(query, ann.Id, imageId, collection, ann.Label.Id, "image", userId, t)
	if err != nil {
		return fmt.Errorf("adding image label annotation record: %v: %w", err, e.ErrInternal)
	}

	return nil
}

func (r AnnotationRepo) findLabelById(labelId l.LabelId) (*l.Label, error) {
	rec := sl.LabelRecord{}
	err := r.Db.Get(&rec,
		"SELECT id,name,description FROM labels WHERE id=$1", labelId)
	if err != nil {
		return nil, fmt.Errorf("fetching label by id %v: %w", labelId, e.ErrInternal)
	}
	return &l.Label{Id: rec.Id, Name: rec.Name, Description: rec.Description}, nil
}

func (r AnnotationRepo) FindImageLabels(
	imageId i.ImageId,
	collection c.CollectionName,
) ([]a.ImageLabel, error) {
	query := `SELECT id,label_id,type,author,touched_at FROM annotations
	WHERE image_id=$1 AND collection_id=(SELECT id FROM collections WHERE name=$2) AND type='image'`

	errCtx := "querying image annotations"
	records := []AnnotationRow{}
	if err := r.Db.Select(&records, query, imageId, collection); err != nil {
		return nil, fmt.Errorf("%v: applying query: %v: %w", errCtx, err, e.ErrInternal)
	}

	imageLabels := []a.ImageLabel{}
	for _, rec := range records {
		label, err := r.findLabelById(rec.LabelId)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", errCtx, err)
		}
		imageLabels = append(
			imageLabels,
			a.ImageLabel{Id: rec.Id, Label: *label, Author: rec.Author, Time: rec.Time},
		)
	}

	return imageLabels, nil
}

func (r AnnotationRepo) RemoveAllAnnotations(imageId i.ImageId, collection string) error {
	_, err := r.Db.Exec(
		"DELETE FROM annotations WHERE image_id=$1 AND collection_id=(SELECT id FROM collections WHERE name=$2)",
		imageId,
		collection,
	)
	if err != nil {
		return fmt.Errorf("deleting annotations: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r AnnotationRepo) RemoveAnnotation(id a.AnnotationId) error {
	_, err := r.Db.Exec("DELETE FROM annotations WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("deleting annotation record: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r AnnotationRepo) RemoveImageLabel(
	imageId i.ImageId,
	collection c.CollectionName,
	labelId l.LabelId,
) error {
	_, err := r.Db.Exec(
		"DELETE FROM annotations WHERE image_id=$1 AND collection_id=(SELECT id FROM collections WHERE name=$2) AND label_id=$3 AND type='image'",
		imageId,
		collection,
		labelId,
	)
	if err != nil {
		return fmt.Errorf("deleting image label: %v: %w", err, e.ErrInternal)
	}
	return nil
}

func (r AnnotationRepo) AddPolygon(
	imageId i.ImageId,
	collection c.CollectionName,
	polygon a.Polygon,
	userId *u.UserId,
	t *time.Time,
) error {
	pointSpecs := []PointSpec{}
	for _, p := range polygon.Points.Coordinates {
		pointSpecs = append(pointSpecs, PointSpec{X: p[0], Y: p[1]})
	}
	coordsBytes, _ := json.Marshal(PolygonSpecs{Points: pointSpecs})
	coordsString := string(coordsBytes)
	query := `INSERT INTO annotations (id, image_id, collection_id, label_id, type, coordinates, author, touched_at)
		VALUES ($1,$2,(SELECT id FROM collections WHERE name=$3),$4,$5,$6,$7,$8)`
	_, err := r.Db.Exec(
		query,
		polygon.Id,
		imageId,
		collection,
		polygon.Label.Id,
		"polygon",
		coordsString,
		userId,
		t,
	)
	if err != nil {
		return fmt.Errorf("inserting polygon: %v: %w", err, e.ErrInternal)
	}

	return nil
}

func (r AnnotationRepo) FindPolygons(
	imageId i.ImageId,
	collection c.CollectionName,
) ([]a.Polygon, error) {
	query := `SELECT id,label_id,type,coordinates,author,touched_at
		FROM annotations
		WHERE image_id=$1 AND collection_id=(SELECT id FROM collections WHERE name=$2) AND type='polygon'`

	errCtx := "querying polygon annotations"
	records := []AnnotationRow{}
	if err := r.Db.Select(&records, query, imageId, collection); err != nil {
		return nil, fmt.Errorf("%v: applying query: %v: %w", errCtx, err, e.ErrInternal)
	}

	polygons := []a.Polygon{}
	for _, rec := range records {
		var specs PolygonSpecs
		err := json.Unmarshal([]byte(rec.Coordinates), &specs)
		if err != nil {
			return nil, fmt.Errorf("%v: unmarshaling polygon specs: %+v: %w: %w",
				errCtx, rec.Coordinates, err, e.ErrInternal)
		}
		label, err := r.findLabelById(rec.LabelId)

		points := a.Points{}
		for _, p := range specs.Points {
			points.Coordinates = append(points.Coordinates, [2]float32{p.X, p.Y})
		}
		polygon := a.NewPolygon(rec.Id, points, *label)

		if rec.Author != nil {
			polygon.Author = rec.Author
		}
		if rec.Time != nil {
			polygon.Time = rec.Time
		}
		polygons = append(polygons, polygon)
	}

	return polygons, nil
}

func (r AnnotationRepo) AddBoundingBox(
	imageId i.ImageId,
	collection c.CollectionName,
	box a.BoundingBox,
	userId *u.UserId,
	t *time.Time,
) error {
	coordsBytes, _ := json.Marshal(
		BoundingBoxSpecs{
			Xc:     box.Xc,
			Yc:     box.Yc,
			Width:  box.Width,
			Height: box.Height,
			Angle:  box.Angle,
		},
	)
	coordsString := string(coordsBytes)
	query := `INSERT INTO annotations (id, image_id, collection_id, label_id, type, coordinates, author, touched_at)
		VALUES ($1,$2,(SELECT id FROM collections WHERE name=$3),$4,$5,$6,$7,$8)`
	_, err := r.Db.Exec(
		query,
		box.Id,
		imageId,
		collection,
		box.Label.Id,
		"bounding_box",
		coordsString,
		userId,
		t,
	)
	if err != nil {
		return fmt.Errorf("inserting bounding box: %v: %w", err, e.ErrInternal)
	}

	return nil
}

func (r AnnotationRepo) FindBoundingBoxes(
	imageId i.ImageId,
	collection c.CollectionName,
) ([]a.BoundingBox, error) {
	query := `SELECT id,label_id,type,coordinates,author,touched_at
		FROM annotations
		WHERE image_id=$1 AND collection_id=(SELECT id FROM collections WHERE name=$2) AND type='bounding_box'`

	errCtx := "querying bounding-box annotations"
	records := []AnnotationRow{}
	if err := r.Db.Select(&records, query, imageId, collection); err != nil {
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

func (r AnnotationRepo) UpdateLabelOfAnnotation(
	id a.AnnotationId,
	labelId l.LabelId,
	userId *u.UserId,
	t *time.Time,
) error {
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

func (r AnnotationRepo) UpdateBoundingBoxCoordinates(
	id a.AnnotationId,
	xc, yc, width, height, angle float32,
) error {
	errCtx := "updating bounding box coordinates"
	if err := a.ValidateBoundingBox(xc, yc, width, height, angle); err != nil {
		return fmt.Errorf("%v: %w", errCtx, err)
	}

	coordsBytes, _ := json.Marshal(
		BoundingBoxSpecs{Xc: xc, Yc: yc, Width: width, Height: height, Angle: angle},
	)
	coordsString := string(coordsBytes)
	query := "UPDATE annotations SET coordinates=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, coordsString, id)
	if err != nil {
		return fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
	}
	return nil
}

func (r AnnotationRepo) UpdateAuthor(id a.AnnotationId, userId *u.UserId) error {
	query := "UPDATE annotations SET author=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, userId, id)
	if err != nil {
		return fmt.Errorf("%w: %w", err, e.ErrInternal)
	}

	return nil
}

func (r AnnotationRepo) UpdateTime(id a.AnnotationId, t *time.Time) error {
	query := "UPDATE annotations SET touched_at=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, t, id)
	if err != nil {
		return fmt.Errorf("%w: %w", err, e.ErrInternal)
	}

	return nil
}

func (r AnnotationRepo) UpdateBoundingBox(
	id a.AnnotationId,
	u a.BoundingBoxUpdatables,
	userId *u.UserId,
	t *time.Time,
) error {
	errCtx := "updating bounding box"
	if err := r.UpdateLabelOfAnnotation(id, u.LabelId, userId, t); err != nil {
		return fmt.Errorf("%v: updating label: %w", errCtx, err)
	}

	if err := r.UpdateBoundingBoxCoordinates(
		id,
		u.Xc,
		u.Yc,
		u.Width,
		u.Height,
		u.Angle,
	); err != nil {
		return fmt.Errorf("%v: updating coordinates: %w", errCtx, err)
	}
	return nil
}

func (r AnnotationRepo) UpdatePolygonPoints(id a.AnnotationId, points a.Points) error {
	errCtx := "updating polygon points"
	pointSpecs := []PointSpec{}
	for _, p := range points.Coordinates {
		pointSpecs = append(pointSpecs, PointSpec{X: p[0], Y: p[1]})
	}
	coordsBytes, _ := json.Marshal(PolygonSpecs{Points: pointSpecs})
	coordsString := string(coordsBytes)
	query := "UPDATE annotations SET coordinates=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, coordsString, id)
	if err != nil {
		return fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrInternal)
	}
	return nil
}

func (r AnnotationRepo) UpdatePolygon(
	id a.AnnotationId,
	u a.PolygonUpdatables,
	userId *u.UserId,
	t *time.Time,
) error {
	errCtx := "updating polygon"
	if err := r.UpdateLabelOfAnnotation(id, u.LabelId, userId, t); err != nil {
		return fmt.Errorf("%v: updating label: %w", errCtx, err)
	}

	if err := r.UpdatePolygonPoints(id, u.Points); err != nil {
		return fmt.Errorf("%v: updating polygon points: %w", errCtx, err)
	}
	return nil
}

func (r AnnotationRepo) GroupOfAnnotation(id a.AnnotationId) (*string, error) {
	var group string
	err := r.Db.Get(
		&group,
		`SELECT name FROM groups WHERE id=(SELECT group_id FROM collections WHERE id=(SELECT collection_id FROM annotations WHERE id=$1))`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetching group of annotation by id %v: %w", id, e.ErrInternal)
	}
	return &group, nil
}

func NewAnnotationRepo(db adb.Querier) AnnotationRepo {
	return AnnotationRepo{Db: db}
}
