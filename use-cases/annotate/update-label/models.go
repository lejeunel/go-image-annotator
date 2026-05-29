package update_label

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Response struct {
	AnnotationId a.AnnotationId
	Label        string
}

type Request struct {
	AnnotationId string
	Label        string
}
