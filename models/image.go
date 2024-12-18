package models

import (
	"github.com/google/uuid"
	"time"
)

type Image struct {
	Id            uuid.UUID `db:"id"`
	Uri           string    `db:"uri"`
	SHA256        string    `db:"sha256"`
	MIMEType      string    `db:"mimetype"`
	Width         int       `db:"width"`
	Height        int       `db:"height"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	Data          []byte
	Annotations   []*Annotation
	BoundingBoxes []*BoundingBox
}
