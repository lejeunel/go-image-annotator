package add_bbox

import (
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Response struct {
	AnnotationId an.AnnotationId
	ImageId      im.ImageId
	Collection   string
	Label        string
	Xc           float32
	Yc           float32
	Width        float32
	Height       float32
}

type Request struct {
	ImageId    im.ImageId
	Collection string
	Label      string
	Xc         float32
	Yc         float32
	Width      float32
	Height     float32
}
