package presenters

import (
	"github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	fetchlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/label/fetch-all"
)

type InitPresenter struct {
	v.View
}

func (p InitPresenter) SuccessReadImage(im im.Image) {
	p.View.DrawImageInfo(view.NewImageInfo(im.Id, im.Collection.Name, im.Specs))
	p.View.DrawImage(view.NewImage(im.Id, im.Reader, im.Collection.Name, im.MIMEType))
	p.View.DrawAnnotationList(view.MakeBoundingBoxes(im.BoundingBoxes),
		view.MakeImageLabels(im.Labels))
}
func (p InitPresenter) SuccessFetchLabels(r fetchlbl.Response) {
	p.View.SetAvailableLabels(r.Labels)
}
func (p InitPresenter) Error(err error) {}
