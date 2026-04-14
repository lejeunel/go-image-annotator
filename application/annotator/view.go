package annotator

import (
	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
)

type AnnotatorView interface {
	DrawScroller(ScrollerButtons)
	Error(error)
	DrawImage(Image)
	DrawImageInfo(ImageInfo)
	AddBox(a.BoundingBox)
	SetAvailableLabels([]string)
	UpdateBox(updbox.Response)
	DeleteAnnotation(del.Response)
}
