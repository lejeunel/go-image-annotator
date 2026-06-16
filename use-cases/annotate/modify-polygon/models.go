package modify_polygon

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Response struct {
	Id a.AnnotationId
}

type Request struct {
	AnnotationId string
	Label        string
	Points       a.Points
}
