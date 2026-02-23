package images

import (
	"database/sql"
	clc "datahub/domain/collections"
	lbl "datahub/domain/labels"
	e "datahub/errors"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type SQLAnnotationRepo struct {
	Db             *sqlx.DB
	ErrorConverter e.DBErrorConverter
}

func NewSQLiteAnnotationRepo(db *sqlx.DB) *SQLAnnotationRepo {

	return &SQLAnnotationRepo{Db: db, ErrorConverter: &e.SQLiteErrorConverter{}}

}

func NewPostgreSQLAnnotationRepo(db *sqlx.DB) *SQLAnnotationRepo {

	return &SQLAnnotationRepo{Db: db, ErrorConverter: &e.PostgreSQLErrorConverter{}}

}

func (r *SQLAnnotationRepo) GetAnnotationIdsOfImage(image *Image) ([]string, error) {
	var annotationsIds []string
	err := r.Db.Select(&annotationsIds, "SELECT id FROM annotations WHERE image_id = $1 AND collection_id=$2 AND shape_type=''",
		image.Id, image.Collection.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrDB
	}
	return annotationsIds, nil
}

func (r *SQLAnnotationRepo) GetAnnotationById(id string) (*Annotation, error) {
	annotation := Annotation{}
	err := r.Db.Get(&annotation, "SELECT id,label_id,image_id,collection_id,created_at,updated_at,author_email FROM annotations WHERE id=$1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNotFound
		}
		return nil, e.ErrDB
	}

	return &annotation, nil
}

func (r *SQLAnnotationRepo) getShapeDataFromId(id string) (string, error) {

	var data string
	err := r.Db.Get(&data, "SELECT shape_data FROM annotations WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrNotFound
		}
		return "", e.ErrDB
	}

	return data, nil

}

func (r *SQLAnnotationRepo) isLabelAlreadyAppliedToImage(label *lbl.Label, image *Image, authorEmail string) (bool, error) {
	var count int64
	query := "SELECT COUNT(*) FROM annotations WHERE label_id=$1 AND image_id=$2 AND collection_id=$3 AND author_email=$4"
	if err := r.Db.QueryRow(query, label.Id, image.Id, image.Collection.Id, authorEmail).Scan(&count); err != nil {
		return false, e.ErrDB
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r *SQLAnnotationRepo) ApplyLabelToImage(label *lbl.Label, image *Image, authorEmail string) error {
	hasDuplicate, err := r.isLabelAlreadyAppliedToImage(label, image, authorEmail)
	if err != nil {
		return err
	}
	if hasDuplicate {
		return e.ErrDuplication
	}

	now := time.Now()
	query := "INSERT INTO annotations (id,image_id,label_id,collection_id,author_email,created_at,updated_at,shape_type,shape_data) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)"

	_, err = r.Db.Exec(query, uuid.New(), image.Id.String(), label.Id, image.Collection.Id, authorEmail, now, now, "", "")
	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLAnnotationRepo) DeleteAllAnnotations(collection *clc.Collection) error {
	_, err := r.Db.Exec("DELETE FROM annotations WHERE collection_id=$1",
		collection.Id)
	if err != nil {
		return e.ErrDB
	}
	return nil
}

func (r *SQLAnnotationRepo) DeleteAnnotation(annotation *Annotation) error {

	_, err := r.Db.Exec("DELETE FROM annotations WHERE id=$1", annotation.Id.String())
	if err != nil {
		return e.ErrDB
	}
	return nil

}

func (r *SQLAnnotationRepo) applyAnnotatedShapeToImage(annotation *AnnotatedShape,
	shapeData string, shapeType string, image *Image) error {

	annotation.Id = uuid.New()

	query := "INSERT INTO annotations (id,image_id,label_id,collection_id,shape_type,shape_data,author_email,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)"
	_, err := r.Db.Exec(query, annotation.Id, image.Id, annotation.LabelId,
		image.Collection.Id, shapeType, shapeData,
		annotation.AuthorEmail, annotation.CreatedAt, annotation.UpdatedAt)

	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLAnnotationRepo) updateAnnotatedShapeToImage(annotation *AnnotatedShape,
	shapeData string) error {

	query := "UPDATE annotations SET shape_data=$1, author_email=$2, updated_at=$3 WHERE id=$4;"
	_, err := r.Db.Exec(query, shapeData,
		annotation.AuthorEmail, time.Now(), annotation.Id)

	if err != nil {
		return e.ErrDB
	}

	return nil
}

func (r *SQLAnnotationRepo) GetBoundingBoxesOfImage(image *Image) ([]*BoundingBox, error) {

	var ids []string
	var bboxes []*BoundingBox

	err := r.Db.Select(&ids, "SELECT id FROM annotations WHERE image_id = $1 AND shape_type='bounding_box' AND collection_id=$2 ORDER BY created_at DESC",
		image.Id, image.Collection.Id)

	if err != nil {
		return nil, e.ErrNotFound
	}

	for _, id := range ids {
		shapeStr, err := r.getShapeDataFromId(id)
		if err != nil {
			return nil, err
		}
		annotation, err := r.GetAnnotationById(id)
		if err != nil {
			return nil, err
		}

		bbox, err := NewBoundingBoxFromJSONShape(shapeStr)
		bbox.Annotation = AnnotatedShape{Annotation: *annotation,
			ShapeData: shapeStr, ShapeType: "bounding_box"}
		bboxes = append(bboxes, bbox)

	}

	return bboxes, nil
}

func (r *SQLAnnotationRepo) UpdateBoundingBox(bbox *BoundingBox, image *Image) error {
	shapeData, err := bbox.MarshalCoordsToJSON()
	if err != nil {
		return err
	}
	return r.updateAnnotatedShapeToImage(&bbox.Annotation, string(shapeData))
}

func (r *SQLAnnotationRepo) ApplyBoundingBox(bbox *BoundingBox, image *Image) error {
	shapeData, err := bbox.MarshalCoordsToJSON()
	if err != nil {
		return err
	}
	return r.applyAnnotatedShapeToImage(&bbox.Annotation, string(shapeData), "bounding_box", image)
}

func (r *SQLAnnotationRepo) UpdateAnnotationLabel(annotationId string, labelId string) error {
	query := "UPDATE annotations SET label_id=$1 WHERE id=$2;"
	_, err := r.Db.Exec(query, labelId, annotationId)

	if err != nil {
		return e.ErrDB
	}
	return nil

}
