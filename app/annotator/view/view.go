package view

type View interface {
	SetScroller(ScrollerButtons)
	Error(error)
	AddBox(BoundingBox)
	AddLabel(ImageLabel)
	SetAvailableLabels([]string)
	SetImageInfo(ImageInfo)
	SetImage(Image)
	SetAnnotations([]BoundingBox, []ImageLabel)
	UpdateBox(BoundingBox)
	UpdateLabel(Annotation)
	DeleteAnnotation(string)
}
