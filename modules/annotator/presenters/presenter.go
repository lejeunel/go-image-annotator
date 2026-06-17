package presenters

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
	fetchlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
)

type Presenter struct {
	v.View
	Colorizer
	usedImageLabels []string
}

func New() Presenter {
	return Presenter{Colorizer: NewCyclicColorizer(Palette)}
}

func (p *Presenter) SetView(view v.View) *Presenter {
	p.View = view
	return p
}

func (p Presenter) SuccessDeleteAnnotation(r del.Response) {
	p.View.DeleteAnnotation(r.Id.String())
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
func (p Presenter) SuccessAddLabel(r addlbl.Response)    {}
func (p Presenter) SuccessAddBox(r addbox.Response)      {}
func (p Presenter) SuccessAddPolygon(r addpoly.Response) {}
func (p Presenter) SuccessUpdateBox(r updbox.Response)   {}
func (p Presenter) SuccessUpdateLabel(r updlbl.Response) {
	p.View.UpdateLabel(v.Annotation{Id: r.AnnotationId.String(), Label: r.Label})
}

func (p Presenter) Error(err error) {}
