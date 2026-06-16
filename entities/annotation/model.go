package annotation

import (
	"fmt"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"time"
)

type Points = [][2]float32

type ImageLabel struct {
	Id     AnnotationId
	Label  lbl.Label
	Author *u.UserId
	Time   *time.Time
}

type Annotation struct {
	Id    AnnotationId
	Label string
}

type BoundingBox struct {
	Id     AnnotationId
	Label  lbl.Label
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
	Angle  float32
	Author *u.UserId
	Time   *time.Time
}

type Polygon struct {
	Id     AnnotationId
	Label  lbl.Label
	Points Points
	Author *u.UserId
	Time   *time.Time
}

type PolygonUpdatables struct {
	LabelId lbl.LabelId
	Points  Points
}

type BoundingBoxResponse struct {
	Label  string
	Xc     float32
	Yc     float32
	Width  float32
	Height float32
}

type BoundingBoxUpdatables struct {
	LabelId lbl.LabelId
	Xc      float32
	Yc      float32
	Width   float32
	Height  float32
	Angle   float32
}

type Option func(*BoundingBox)

func WithAngle(ang float32) Option {
	return func(c *BoundingBox) {
		c.Angle = ang
	}
}

func NewBoundingBox(id AnnotationId, xc float32, yc float32,
	width float32, height float32, label lbl.Label, opts ...Option,
) BoundingBox {
	b := &BoundingBox{Id: id, Xc: xc, Yc: yc, Width: width, Height: height, Label: label}
	for _, opt := range opts {
		opt(b)
	}
	return *b
}

func ValidateBoundingBox(xc float32, yc float32, width float32, height float32, angle float32) error {
	errCtx := "validating bounding box"
	if width <= 0 {
		return fmt.Errorf("%v: checking whether width (%v) <= 0: %w", errCtx, width, e.ErrValidation)
	}
	return nil

}

func NewImageLabel(label lbl.Label) ImageLabel {
	return ImageLabel{Id: NewAnnotationId(), Label: label}
}

func NewPolygon(id AnnotationId, points [][2]float32, label lbl.Label) Polygon {
	return Polygon{Id: id, Points: points, Label: label}
}
