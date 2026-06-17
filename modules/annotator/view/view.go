package view

type View interface {
	SetScroller(ScrollerButtons)
	Error(error)
	SetAvailableLabels([]string)
	SetAvailableImageLabels([]string)
	SetImageInfo(ImageInfo)
	SetImage(Image)
	SetAnnotations([]BoundingBox, []Polygon, []ImageLabel)
}
