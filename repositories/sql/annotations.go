package repositories

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	"time"
)

type SQLAnnotationRepo struct {
	Db *sqlx.DB
}

func NewSQLLabelRepo(db *sqlx.DB) *SQLAnnotationRepo {

	return &SQLAnnotationRepo{Db: db}

}

func (r *SQLAnnotationRepo) CreateLabel(ctx context.Context, label *m.Label) (*m.Label, error) {
	now := time.Now().String()
	query := "INSERT INTO labels (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := r.Db.Exec(query, label.Id, label.Name, label.Description, now,
		now)

	if err != nil {
		return nil, err
	}

	return label, nil
}

func (r *SQLAnnotationRepo) NumImagesWithLabel(ctx context.Context, label *m.Label) (int, error) {
	var nImages int
	err := r.Db.QueryRow("SELECT COUNT(*) FROM image_label_assoc WHERE label_id = ?",
		label.Id).Scan(&nImages)
	if err != nil {
		return 0, err
	}

	return nImages, nil

}

func (r *SQLAnnotationRepo) DeleteLabel(ctx context.Context, label *m.Label) error {
	_, err := r.Db.Exec("DELETE FROM labels WHERE id=?", label.Id.String())
	return err
}

func (r *SQLAnnotationRepo) DeletePolygon(ctx context.Context, polygon *m.Polygon) error {
	_, err := r.Db.Exec("DELETE FROM polygons WHERE id=?", polygon.Id.String())
	return err
}

func (r *SQLAnnotationRepo) GetOneLabel(ctx context.Context, id string) (*m.Label, error) {
	label := m.Label{}
	err := r.Db.Get(&label, "SELECT id,name,description FROM labels WHERE id=?", id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "label", Criteria: "id", Value: id, Err: err}
	}

	return &label, nil
}

func (r *SQLAnnotationRepo) getOnePolygon(ctx context.Context, id string) (*m.Polygon, error) {
	polygon := m.Polygon{}
	err := r.Db.Get(&polygon, "SELECT id,type_,min_x,min_y,max_x,max_y,created_at,updated_at FROM polygons WHERE id=?", id)
	if err != nil {
		return nil, &e.ErrNotFound{Entity: "polygon", Criteria: "id", Value: id, Err: err}
	}

	var pointsStr string
	_ = r.Db.Get(&pointsStr, "SELECT points FROM polygons WHERE id=?", id)
	var points [][]int
	if err = json.Unmarshal([]byte(pointsStr), &points); err != nil {
		return nil, err
	}

	polygon.Points = points

	var labelId string
	err = r.Db.Get(&labelId, "SELECT label_id FROM polygons WHERE id=?", id)
	if err != nil {
		return nil, &e.ErrNotFound{Entity: "label", Criteria: "id", Value: id, Err: err}
	}

	label, err := r.GetOneLabel(ctx, labelId)
	if err != nil {
		return nil, err
	}
	polygon.Label = label

	return &polygon, nil
}

func (r *SQLAnnotationRepo) GetLabelsOfImage(ctx context.Context, image *m.Image) ([]*m.Label, error) {
	var labelIds []string
	var labels []*m.Label

	err := r.Db.Select(&labelIds, "SELECT label_id FROM image_label_assoc WHERE image_id = ?", image.Id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "image", Criteria: "id", Value: image.Id.String(), Err: err}
	}

	for _, id := range labelIds {
		l, err := r.GetOneLabel(ctx, id)
		if err != nil {
			return nil, err
		}
		labels = append(labels, l)

	}

	return labels, nil

}

func (r *SQLAnnotationRepo) ApplyLabelToImage(ctx context.Context, label *m.Label, image *m.Image) error {
	now := time.Now().String()
	query := "INSERT INTO image_label_assoc (image_id, label_id, created_at) VALUES (?, ?, ?)"
	_, err := r.Db.Exec(query, image.Id, label.Id, now)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLAnnotationRepo) ApplyPolygonToImage(ctx context.Context, polygon *m.Polygon, image *m.Image) error {
	now := time.Now().String()

	points, err := json.Marshal(polygon.Points)
	if err != nil {
		return err
	}

	query := "INSERT INTO polygons (id,image_id,type_,min_x,min_y,max_x,max_y,points,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?)"
	_, err = r.Db.Exec(query, polygon.Id, image.Id,
		polygon.Type, polygon.MinX, polygon.MinY, polygon.MaxX, polygon.MaxY, string(points), now, now)

	if err != nil {
		return err
	}

	if polygon.Label != nil {
		query := "UPDATE polygons SET label_id=? WHERE id=?"
		_, err = r.Db.Exec(query, polygon.Label.Id, polygon.Id)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLAnnotationRepo) GetPolygonsOfImage(ctx context.Context, image *m.Image) ([]*m.Polygon, error) {

	var polygonIds []string
	var polygons []*m.Polygon

	err := r.Db.Select(&polygonIds, "SELECT id FROM polygons WHERE image_id = ?", image.Id)

	if err != nil {
		return nil, &e.ErrNotFound{Entity: "image", Criteria: "id", Value: image.Id.String(), Err: err}
	}

	for _, id := range polygonIds {
		p, err := r.getOnePolygon(ctx, id)
		if err != nil {
			return nil, err
		}
		polygons = append(polygons, p)

	}

	return polygons, nil
}

func (r *SQLAnnotationRepo) Nums() (int64, error) {

	return 0, nil

}

func (r *SQLAnnotationRepo) Slice(offset, length int, data interface{}) error {

	return nil
}