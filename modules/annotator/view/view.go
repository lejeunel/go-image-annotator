package view

import (
	m "github.com/lejeunel/go-image-annotator/entities/meta"
)

type View interface {
	SetScroller(ScrollerButtons)
	Error(error)
	SetAvailableLabels([]string)
	SetAvailableImageLabels([]string)
	SetImageInfo(ImageInfo)
	SetImage(Image)
	SetAnnotations([]BoundingBox, []Polygon, []ImageLabel)
	SetMetaData([]m.MetaData)
}
