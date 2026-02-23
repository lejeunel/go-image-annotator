package images

import (
	lbl "datahub/domain/labels"
	e "datahub/errors"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type BBoxCoords struct {
	Xc     float64 `json:"xc,string"`
	Yc     float64 `json:"yc,string"`
	Height float64 `json:"height,string"`
	Width  float64 `json:"width,string"`
	Angle  float64 `json:"angle,string"`
}

type BoundingBox struct {
	Annotation AnnotatedShape
	Coords     BBoxCoords
}

func (b BoundingBox) ToMap() map[string]any {
	return map[string]any{
		"id":     b.Annotation.Id,
		"xc":     fmt.Sprintf("%f", b.Coords.Xc),
		"yc":     fmt.Sprintf("%f", b.Coords.Yc),
		"height": fmt.Sprintf("%f", b.Coords.Height),
		"width":  fmt.Sprintf("%f", b.Coords.Width),
		"angle":  fmt.Sprintf("%f", b.Coords.Angle),
	}
}

func (b BoundingBox) MarshalCoordsToJSON() ([]byte, error) {

	return json.Marshal(b.ToMap())
}

func (b *BoundingBox) Validate() error {
	if err := validation.ValidateStruct(&b.Coords,
		validation.Field(&b.Coords.Xc, validation.Min(0.)),
		validation.Field(&b.Coords.Yc, validation.Min(0.))); err != nil {
		return fmt.Errorf("all coordinates are required, and must be >= 0: %w",
			e.ErrValidation)
	}
	return nil
}

func (b *BoundingBox) Annotate(label *lbl.Label) {
	b.Annotation.LabelId = label.Id
	b.Annotation.Label = label
}

func NewBoundingBoxFromJSONShape(str string) (*BoundingBox, error) {
	var coords BBoxCoords

	if err := json.Unmarshal([]byte(str), &coords); err != nil {
		return nil, err
	}
	return NewBoundingBox(coords.Xc, coords.Yc, coords.Height, coords.Width)
}

func NewBoundingBox(xc, yc, h, w float64) (*BoundingBox, error) {
	coords := BBoxCoords{Xc: xc, Yc: yc, Height: h, Width: w}
	annotation := AnnotatedShape{Annotation: Annotation{Id: uuid.New()}, ShapeType: "bounding_box"}
	bbox := &BoundingBox{Coords: coords, Annotation: annotation}
	if err := bbox.Validate(); err != nil {
		return nil, fmt.Errorf("upserting bounding box validation: %w", err)
	}
	return bbox, nil
}

func NewBoundingBoxFromMinMax(xmin, ymin, xmax, ymax float64) *BoundingBox {
	height := ymax - ymin
	width := xmax - xmin
	return &BoundingBox{Coords: BBoxCoords{Xc: xmin + width/2, Yc: ymin + height/2,
		Height: height, Width: width}}
}
