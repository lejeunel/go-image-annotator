package assign_label

import (
	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Response struct {
	AnnotationId an.AnnotationId
	ImageId      im.ImageId
	Collection   string
	Label        string
}

type Request struct {
	ImageId    im.ImageId
	Collection string
	Label      string
}
