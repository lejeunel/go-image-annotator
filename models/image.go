package models

import (
	"github.com/google/uuid"
)

type Image struct {
	Id          uuid.UUID `db:"id"`
	Uri         string    `db:"uri"`
	SHA256      string    `db:"sha256"`
	MIMEType    string    `db:"mimetype"`
	Width       int       `db:"width"`
	Height      int       `db:"height"`
	CreatedAt   string    `db:"created_at"`
	UpdatedAt   string    `db:"updated_at"`
	Data        []byte
	Annotations []*ImageAnnotation
	Polygons    []*Polygon
}

// func (im Image) Validate() error {
// 	return validation.ValidateStruct(&a,
// 		validation.Field(&a.FirstName, validation.Required),
// 		validation.Field(&a.LastName, validation.Required),
// 		validation.Field(&a.DateOfBirth, validation.Date("2006-01-02")),
// 	)
// }
