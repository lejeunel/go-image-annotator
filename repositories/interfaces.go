package repositories

import (
	c "context"
	pag "github.com/vcraescu/go-paginator/v2"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
)

type ImageRepo interface {
	Create(c.Context, *m.Image) (*m.Image, error)
	Delete(c.Context, *m.Image) error
	GetOne(c.Context, string) (*m.Image, error)
	Paginate(pageSize int, filters *g.ImageFilterArgs) pag.Paginator
}

type AnnotationRepo interface {
	CreateLabel(c.Context, *m.Label) (*m.Label, error)
	DeleteLabel(c.Context, *m.Label) error
	GetOneLabel(c.Context, string) (*m.Label, error)

	GetLabelsOfImage(c.Context, *m.Image) ([]*m.Label, error)
	ApplyLabelToImage(c.Context, *m.Label, *m.Image) error
	NumImagesWithLabel(c.Context, *m.Label) (int, error)

	ApplyPolygonToImage(c.Context, *m.Polygon, *m.Image) error
	GetPolygonsOfImage(c.Context, *m.Image) ([]*m.Polygon, error)
	DeletePolygon(c.Context, *m.Polygon) error

	Nums() (int64, error)
	Slice(offset, length int, data interface{}) error
}

type SetRepo interface {
	Create(c.Context, *m.Set) (*m.Set, error)
	GetOne(c.Context, string) (*m.Set, error)
	AssignImageToSet(c.Context, *m.Image, *m.Set) error
	Nums() (int64, error)
	Slice(offset, length int, data interface{}) error
}
