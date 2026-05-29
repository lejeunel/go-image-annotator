package image

import (
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
)

type Response struct {
	Id            ImageId
	Collection    string
	Labels        []*a.ImageLabel
	BoundingBoxes []*a.BoundingBox
}
