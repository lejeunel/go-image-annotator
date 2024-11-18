package models

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type Polygon struct {
	Id        uuid.UUID `db:"id"`
	Type      string    `db:"type_"`
	MinX      int       `db:"min_x"`
	MinY      int       `db:"min_y"`
	MaxX      int       `db:"max_x"`
	MaxY      int       `db:"max_y"`
	CreatedAt string    `db:"created_at"`
	UpdatedAt string    `db:"updated_at"`
	Points    [][]int
	Label     *Label
}

func (p *Polygon) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Type, validation.Required),
		validation.Field(&p.Type, validation.In("rectangle", "polygon")),
		validation.Field(&p.MinX, validation.Min(0)),
		validation.Field(&p.MinY, validation.Min(0)))
}

func validateBoundingBoxPoints(bbox *Polygon) error {
	if bbox.Points[0][0] > bbox.Points[1][0] {
		return errors.New(fmt.Sprintf("BoundingBox should have x0 > x1, got %v<=%v", bbox.Points[0][0],
			bbox.Points[1][0]))
	}

	if bbox.Points[0][1] > bbox.Points[1][1] {
		return errors.New(fmt.Sprintf("BoundingBox should have y0 > y1, got %v<=%v", bbox.Points[0][1],
			bbox.Points[1][1]))
	}

	return nil

}

func NewBoundingBox(x0, y0, x1, y1 int) (*Polygon, error) {
	p := &Polygon{Id: uuid.New(), Type: "rectangle",
		MinX: x0, MinY: y0, MaxX: x1, MaxY: y1,
		Points: [][]int{[]int{x0, y0}, []int{x1, y1}},
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	if err := validateBoundingBoxPoints(p); err != nil {
		return nil, err
	}

	return p, nil

}
