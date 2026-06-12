package modify_bbox

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Response struct {
	Id a.AnnotationId
}

type Request struct {
	AnnotationId string
	Label        string
	Xc           float32
	Yc           float32
	Width        float32
	Height       float32
	Angle        float32
}
