package images

import (
	clc "datahub/domain/collections"
	lbl "datahub/domain/labels"
	"github.com/google/uuid"
	"time"
)

type Annotation struct {
	Id           uuid.UUID        `db:"id"`
	ImageId      ImageId          `db:"image_id"`
	LabelId      lbl.LabelId      `db:"label_id"`
	CollectionId clc.CollectionId `db:"collection_id"`
	AuthorEmail  string           `db:"author_email"`
	CreatedAt    time.Time        `db:"created_at"`
	UpdatedAt    time.Time        `db:"updated_at"`
	Label        *lbl.Label
}

type AnnotatedShape struct {
	Annotation
	ShapeData string `db:"shape_data"`
	ShapeType string `db:"shape_type"`
}
