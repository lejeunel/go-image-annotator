package view

import (
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
)

type View interface {
	DrawScroller(ScrollerButtons)
	Error(error)
	DrawImage(Image)
	DrawImageInfo(ImageInfo)
	DrawAnnotationList([]*BoundingBox, []*ImageLabel)
	AddBox(BoundingBox)
	SetAvailableLabels([]string)
	UpdateBox(updbox.Response)
	DeleteAnnotation(del.Response)
}
