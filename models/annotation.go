package models

import (
	"github.com/google/uuid"
)

type Annotation struct {
	Id          uuid.UUID `db:"id"`
	ImageId     uuid.UUID `db:"image_id"`
	LabelId     uuid.UUID `db:"label_id"`
	AuthorEmail string    `db:"author_email"`
	CreatedAt   string    `db:"created_at"`
	UpdatedAt   string    `db:"updated_at"`
	Label       *Label
}

type AnnotatedShape struct {
	Annotation
	ShapeData string `db:"shape_data"`
	ShapeType string `db:"shape_type"`
}
