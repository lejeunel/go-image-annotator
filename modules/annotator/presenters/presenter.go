package presenters

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	fetchlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
)

type Presenter struct {
	v.View
	Colorizer
}

func NewAnnotationPagePresenter(view v.View) Presenter {
	return Presenter{View: view, Colorizer: NewCyclicColorizer(Palette)}
}

func (p *Presenter) SetView(view v.View) *Presenter {
	p.View = view
	return p
}
func (p Presenter) SuccessInitScroller(s scr.ScrollerState) {
	p.View.SetScroller(MakeScrollerButtons(s))
}

func (p Presenter) SuccessReadImage(im im.Image) {
	p.View.SetImageInfo(v.NewImageInfo(im.Id, im.Collection.Name, im.Specs))
	p.View.SetImage(v.NewImage(im.Id, im.Reader, im.Collection.Name, im.MIMEType))
	p.View.SetAnnotations(MakeBoundingBoxes(im.BoundingBoxes, p.Colorizer),
		MakePolygons(im.Polygons, p.Colorizer),
		MakeImageLabels(im.Labels))
}
func (p Presenter) SuccessFetchLabels(r fetchlbl.Response) {
	p.View.SetAvailableLabels(r.Labels)
}

func (p Presenter) Error(err error) {
	p.View.Error(err)
}
