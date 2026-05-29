package modify_bbox

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Response struct {
	AnnotationId a.AnnotationId
	Label        string
	Xc           float32
	Yc           float32
	Width        float32
	Height       float32
}

type Request struct {
	AnnotationId a.AnnotationId
	Label        string
	Xc           float32
	Yc           float32
	Width        float32
	Height       float32
}
