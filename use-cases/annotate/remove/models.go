package remove

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Response struct {
	Id a.AnnotationId
}

type Request struct {
	Id a.AnnotationId
}
