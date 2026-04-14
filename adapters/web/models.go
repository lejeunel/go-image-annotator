package web

import (
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
)

type Bounds struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

type BoxGeometry struct {
	XTopLeft float64 `json:"x"`
	YTopLeft float64 `json:"y"`
	W        float64 `json:"w"`
	H        float64 `json:"h"`
	Bounds   Bounds  `json:"bounds"`
}

type BoxSelector struct {
	Type     string      `json:"type"`
	Geometry BoxGeometry `json:"geometry"`
}

type BoxAnnotation struct {
	AnnotationId string      `json:"annotation"`
	Selector     BoxSelector `json:"selector"`
}

type Properties struct {
	Color string `json:"color"`
	Label string `json:"label"`
}

type BoxRequest struct {
	ImageId    string        `json:"image_id"`
	Collection string        `json:"collection"`
	Annotation BoxAnnotation `json:"annotation"`
	Label      string        `json:"label"`
}

type Request struct {
	ImageId    im.ImageId
	Collection string
}
