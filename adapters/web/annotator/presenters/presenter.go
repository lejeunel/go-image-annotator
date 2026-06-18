package presenters

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	fetchlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
)

type AnnotationPagePresenter struct {
	v.View
	Colorizer
}

func NewAnnotationPagePresenter(colorizer Colorizer) AnnotationPagePresenter {
	return AnnotationPagePresenter{Colorizer: colorizer}
}

func (p *AnnotationPagePresenter) SetView(view v.View) *AnnotationPagePresenter {
	p.View = view
	return p
}
func (p AnnotationPagePresenter) SuccessInitScroller(s scr.ScrollerState) {
	p.View.SetScroller(MakeScrollerButtons(s))
}
func (p AnnotationPagePresenter) SuccessReadImage(im im.Image) {
	p.View.SetImageInfo(v.NewImageInfo(im.Id, im.Collection.Name, im.Specs))
	p.View.SetImage(v.NewImage(im.Id, im.Reader, im.Collection.Name, im.MIMEType))
	p.View.SetAnnotations(MakeBoundingBoxes(im.BoundingBoxes, p.Colorizer),
		MakePolygons(im.Polygons, p.Colorizer),
		MakeImageLabels(im.Labels))
}
func (p AnnotationPagePresenter) SuccessFetchLabels(r fetchlbl.Response) {
	p.View.SetAvailableLabels(r.Labels)
}
func (p AnnotationPagePresenter) Error(err error) {
	p.View.Error(err)
}
