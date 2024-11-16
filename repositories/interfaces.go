package repositories

import (
	c "context"
	m "go-image-annotator/models"
)

type ImageRepo interface {
	Create(c.Context, *m.Image) (*m.Image, error)
	GetOne(c.Context, string) (*m.Image, error)
	Nums() (int64, error)
	Slice(offset, length int, data interface{}) error
}

type LabelRepo interface {
	Create(c.Context, *m.Label) (*m.Label, error)
	GetOne(c.Context, string) (*m.Label, error)
	Nums() (int64, error)
	Slice(offset, length int, data interface{}) error
}
