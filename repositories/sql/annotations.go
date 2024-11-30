package repositories

import (
	c "context"
	"github.com/google/uuid"
	pag "github.com/vcraescu/go-paginator/v2"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLAnnotationRepo struct {
	Db *sqlx.DB
}

type PaginableSQLAnnotationRepo struct {
	Repo *SQLAnnotationRepo
}

func NewSQLLabelRepo(db *sqlx.DB) *SQLAnnotationRepo {

	return &SQLAnnotationRepo{Db: db}

}

func (r *SQLAnnotationRepo) CreateLabel(ctx c.Context, label *m.Label) (*m.Label, error) {
	now := time.Now()

	query := "INSERT INTO labels (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := r.Db.Exec(query, label.Id, label.Name, label.Description, now,
		now)

	if err != nil {
		return nil, err
	}

	return label, nil
}

func (r *SQLAnnotationRepo) DeleteLabel(ctx c.Context, label *m.Label) error {
	_, err := r.Db.Exec("DELETE FROM labels WHERE id=?", label.Id.String())
	return err
}

func (r *SQLAnnotationRepo) GetOneLabel(ctx c.Context, id string) (*m.Label, error) {
	label := m.Label{}
	err := r.Db.Get(&label, "SELECT id,name,description FROM labels WHERE id=?", id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "label", Criteria: "id", Value: id, Err: err}
	}

	return &label, nil
}

func (r *SQLAnnotationRepo) getAnnotationById(ctx c.Context, id string) (*m.Annotation, error) {
	annotation := m.Annotation{}
	err := r.Db.Get(&annotation, "SELECT id,label_id,image_id,author_email FROM annotations WHERE id=?", id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "annotation", Criteria: "id", Value: id, Err: err}
	}

	return &annotation, nil
}

func (r *SQLAnnotationRepo) getShapeFromId(ctx c.Context, id string) (string, error) {

	var data string
	err := r.Db.Get(&data, "SELECT shape_data FROM annotations WHERE id=?", id)
	if err != nil {
		return "", err
	}

	return data, nil

}

func (r *SQLAnnotationRepo) GetAnnotationsOfImage(ctx c.Context, image *m.Image, collection *m.Collection) ([]*m.Annotation, error) {
	var annotationsIds []string
	var annotations []*m.Annotation

	err := r.Db.Select(&annotationsIds, "SELECT id FROM annotations WHERE image_id = ? AND collection_id=? AND shape_type=''",
		image.Id, collection.Id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "image", Criteria: "id", Value: image.Id.String(), Err: err}
	}

	for _, id := range annotationsIds {
		a, err := r.getAnnotationById(ctx, id)
		if err != nil {
			return nil, err
		}
		if a.LabelId != uuid.Nil {
			label, err := r.GetOneLabel(ctx, a.LabelId.String())
			if err != nil {
				return nil, err
			}
			a.Label = label

		}
		annotations = append(annotations, a)

	}

	return annotations, nil

}

func (r *SQLAnnotationRepo) ApplyLabelToImage(ctx c.Context, label *m.Label, image *m.Image, collection *m.Collection, authorEmail string) error {
	now := time.Now()
	query := "INSERT INTO annotations (id,image_id,label_id,collection_id,author_email,created_at,shape_type,shape_data) VALUES (?,?,?,?,?,?,?,?)"

	_, err := r.Db.Exec(query, uuid.New(), image.Id, label.Id, collection.Id, authorEmail, now, "", "")

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLAnnotationRepo) DeleteAnnotation(ctx c.Context, annotation *m.Annotation) error {

	_, err := r.Db.Exec("DELETE FROM annotations WHERE id=?", annotation.Id.String())
	if err != nil {
		return err
	}
	return nil

}

func (r *SQLAnnotationRepo) applyAnnotatedShapeToImage(ctx c.Context, annotation *m.AnnotatedShape,
	shapeData string, shapeType string, image *m.Image) error {

	annotation.Id = uuid.New()

	query := "INSERT INTO annotations (id,image_id,label_id,shape_type,shape_data,author_email,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?)"
	_, err := r.Db.Exec(query, annotation.Id, image.Id, annotation.LabelId,
		shapeType, shapeData,
		annotation.AuthorEmail, annotation.CreatedAt, annotation.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLAnnotationRepo) getAnnotatedShapeFromId(ctx c.Context, id string) (*m.AnnotatedShape, error) {
	shape := m.AnnotatedShape{}

	err := r.Db.Get(&shape, "SELECT id,image_id,label_id,author_email,created_at,updated_at FROM annotations WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	label, err := r.GetOneLabel(ctx, shape.LabelId.String())
	if err != nil {
		return nil, err
	}
	shape.Label = label

	return &shape, nil

}

func (r *SQLAnnotationRepo) GetBoundingBoxesOfImage(ctx c.Context, image *m.Image) ([]*m.BoundingBox, error) {

	var ids []string
	var bboxes []*m.BoundingBox

	err := r.Db.Select(&ids, "SELECT id FROM annotations WHERE image_id = ? AND shape_type=\"bounding_box\"", image.Id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "image", Criteria: "id", Value: image.Id.String(), Err: err}
	}

	for _, id := range ids {
		shapeStr, err := r.getShapeFromId(ctx, id)
		if err != nil {
			return nil, err
		}
		annotation, err := r.getAnnotatedShapeFromId(ctx, id)
		if err != nil {
			return nil, err
		}

		bbox, err := m.NewBoundingBoxFromJSONShape(shapeStr)
		if err != nil {
			return nil, err
		}
		bbox.AnnotatedShape = *annotation
		bboxes = append(bboxes, bbox)

	}

	return bboxes, nil
}

func (r *SQLAnnotationRepo) ApplyBoundingBoxToImage(ctx c.Context, bbox *m.BoundingBox, image *m.Image) error {
	shapeData, err := bbox.MarshalCoordsToJSON()
	if err != nil {
		return err
	}
	return r.applyAnnotatedShapeToImage(ctx, &bbox.AnnotatedShape, string(shapeData), "bounding_box", image)
}

func (r *SQLAnnotationRepo) Paginate(pageSize int) pag.Paginator {
	paginable := &PaginableSQLAnnotationRepo{Repo: r}
	return pag.New(paginable, pageSize)
}

func (r *PaginableSQLAnnotationRepo) Nums() (int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM labels "
	err := r.Repo.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func (r *PaginableSQLAnnotationRepo) Slice(offset, length int, data interface{}) error {

	baseQuery := "SELECT id,name,description FROM labels "
	query := baseQuery + " LIMIT $1 OFFSET $2"
	rows, err := r.Repo.Db.Queryx(query, length, offset)

	if err != nil {
		return err
	}

	s := data.(*[]m.Label)

	for rows.Next() {
		var l m.Label
		err := rows.StructScan(&l)

		if err != nil {
			return err
		}

		*s = append(*s, l)
	}

	return nil
}
