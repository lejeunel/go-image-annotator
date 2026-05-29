package annotation

import (
	"fmt"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type AnnotationKind int

const (
	BoundingBoxAnnotation AnnotationKind = iota
	ImageLabelAnnotation
)

type ImageLabel struct {
	Id    AnnotationId
	Label lbl.Label
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
}

func NewBoundingBox(id AnnotationId, xc float32, yc float32, width float32, height float32, label lbl.Label) *BoundingBox {
	return &BoundingBox{Id: id, Xc: xc, Yc: yc, Width: width, Height: height, Label: label}
}

func ValidateBoundingBox(xc float32, yc float32, width float32, height float32) error {
	errCtx := "validating bounding box"
	if width <= 0 {
		return fmt.Errorf("%v: checking whether width (%v) <= 0: %w", errCtx, width, e.ErrValidation)
	}
	return nil

}

func NewImageLabel(label lbl.Label) *ImageLabel {
	return &ImageLabel{Id: NewAnnotationId(), Label: label}
}
