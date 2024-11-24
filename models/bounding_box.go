package models

import (
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BoundingBox struct {
	AnnotatedShape
	Xc     float64 `json:"xc,string"`
	Yc     float64 `json:"yc,string"`
	Height float64 `json:"height,string"`
	Width  float64 `json:"width,string"`
	Angle  float64 `json:"angle,string"`
}

func (b BoundingBox) MarshalCoordsToJSON() ([]byte, error) {
	m := map[string]string{
		"xc":     fmt.Sprintf("%f", b.Xc),
		"yc":     fmt.Sprintf("%f", b.Yc),
		"height": fmt.Sprintf("%f", b.Height),
		"width":  fmt.Sprintf("%f", b.Width),
		"angle":  fmt.Sprintf("%f", b.Angle),
	}

	return json.Marshal(m)
}

func (b *BoundingBox) Validate(image *Image) error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Xc, validation.Required),
		validation.Field(&b.Xc, validation.Min(0.)),
		validation.Field(&b.Yc, validation.Required),
		validation.Field(&b.Yc, validation.Min(0.)))
}

func (b *BoundingBox) Annotate(label *Label) {
	b.LabelId = label.Id
	b.Label = label
}

func NewBoundingBoxFromJSONShape(str string) (*BoundingBox, error) {
	var bbox BoundingBox

	if err := json.Unmarshal([]byte(str), &bbox); err != nil {
		return nil, err
	}
	return &bbox, nil
}

func NewBoundingBox(xc, yc, h, w float64) *BoundingBox {
	return &BoundingBox{Xc: xc, Yc: yc, Height: h, Width: w}
}
