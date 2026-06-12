package add_bbox

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
	Xc         float32
	Yc         float32
	Width      float32
	Height     float32
	Angle      float32
}
