package annotator

import (
	"context"
	p "github.com/lejeunel/go-image-annotator/modules/annotator/presenters"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
	imread "github.com/lejeunel/go-image-annotator/use-cases/image/read"
	fetchlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
)

type Annotator struct {
	scroller        scr.Interface
	imageReader     imread.Interface
	boxAdder        addbox.Interface
	polygonAdder    addpoly.Interface
	imageLabelAdder addlbl.Interface
	boxUpdater      updbox.Interface
	deleter         del.Interface
	labelFetcher    fetchlbl.Interface
	labelUpdater    updlbl.Interface
	presenter       p.Presenter
}

func (a *Annotator) DeleteAnnotation(ctx context.Context, r del.Request, view v.View) {
	a.deleter.Execute(ctx, r, a.presenter.SetView(view))
}
func (a *Annotator) UpdateLabel(ctx context.Context, r updlbl.Request, view v.View) {
	a.labelUpdater.Execute(ctx, r, a.presenter.SetView(view))
}
func (a *Annotator) UpdateBox(ctx context.Context, r updbox.Request, view v.View) {
	a.boxUpdater.Execute(ctx, r, a.presenter.SetView(view))
}
func (a *Annotator) AddBox(ctx context.Context, r addbox.Request, view v.View) {
	a.boxAdder.Execute(ctx, r, a.presenter.SetView(view))
}
func (a *Annotator) AddPolygon(ctx context.Context, r addpoly.Request, view v.View) {
	a.polygonAdder.Execute(ctx, r, a.presenter.SetView(view))
}
func (a *Annotator) AddLabel(ctx context.Context, r addlbl.Request, view v.View) {
	a.imageLabelAdder.Execute(ctx, r, a.presenter.SetView(view))
}
func (a *Annotator) Init(ctx context.Context, imageId string, collection string, view v.View) {
	scrollerState, err := a.scroller.Init(imageId, scr.WithCollection(collection))
	if err != nil {
		view.Error(err)
		return
	}
	view.SetScroller(p.MakeScrollerButtons(*scrollerState))

	a.presenter.SetView(view)
	a.imageReader.Execute(imread.Request{ImageId: imageId, Collection: collection},
		a.presenter)
	a.labelFetcher.Execute(ctx, a.presenter)
}

func NewAnnotator(scroller scr.Interface, imageMetaReader imread.Interface,
	boxAdder addbox.Interface, boxUpdater updbox.Interface, polygonAdder addpoly.Interface, annotationDeleter del.Interface,
	labelFetcher fetchlbl.Interface, labelUpdater updlbl.Interface, imageLabelAdder addlbl.Interface) Annotator {
	return Annotator{
		scroller:        scroller,
		imageReader:     imageMetaReader,
		boxAdder:        boxAdder,
		boxUpdater:      boxUpdater,
		polygonAdder:    polygonAdder,
		deleter:         annotationDeleter,
		labelFetcher:    labelFetcher,
		labelUpdater:    labelUpdater,
		imageLabelAdder: imageLabelAdder,
		presenter:       p.New(),
	}
}
