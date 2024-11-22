package models

import (
	"github.com/google/uuid"
)

type ImageAnnotation struct {
	Id          uuid.UUID `db:"id"`
	LabelId     string    `db:"label_id"`
	ImageId     string    `db:"image_id"`
	AuthorEmail string    `db:"author_email"`
	CreatedAt   string    `db:"created_at"`
}
