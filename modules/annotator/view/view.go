package view

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
)

type View interface {
	SetQuery(im.FilterStr, im.OrderStr)
	SetScroller(ScrollerButtons)
	Error(error)
	SetAvailableLabels([]string)
	SetAvailableImageLabels([]string)
	SetImageInfo(ImageInfo)
	SetImage(Image)
	SetAnnotations([]BoundingBox, []Polygon, []ImageLabel)
	SetMetaData([]m.MetaData)
}
