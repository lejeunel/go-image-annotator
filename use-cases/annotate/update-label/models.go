package update_label

import (
	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
)

type Response struct {
}

type Request struct {
	AnnotationId a.AnnotationId
	Label        string
}
