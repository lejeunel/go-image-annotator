package add_polygon

import (
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Response struct {
	Id an.AnnotationId
}

type Request struct {
	ImageId    string
	Collection string
	Label      string
	Points     an.Points
}
