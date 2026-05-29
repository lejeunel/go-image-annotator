package view

type View interface {
	DrawScroller(ScrollerButtons)
	Error(error)
	DrawImage(Image)
	DrawImageInfo(ImageInfo)
	DrawAnnotationList([]*BoundingBox, []*ImageLabel)
	AddBox(BoundingBox)
	AddLabel(ImageLabel)
	SetAvailableLabels([]string)
	UpdateBox(BoundingBox)
	UpdateLabel(Annotation)
	DeleteAnnotation(string)
}
