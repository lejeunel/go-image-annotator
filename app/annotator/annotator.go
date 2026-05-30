package annotator

import (
	p "github.com/lejeunel/go-image-annotator/app/annotator/presenters"
	scr "github.com/lejeunel/go-image-annotator/app/annotator/scroller"
	v "github.com/lejeunel/go-image-annotator/app/annotator/view"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
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
	imageLabelAdder addlbl.Interface
	boxUpdater      updbox.Interface
	deleter         del.Interface
	labelFetcher    fetchlbl.Interface
	labelUpdater    updlbl.Interface
	presenter       p.Presenter
}

func (a *Annotator) DeleteAnnotation(r del.Request, view v.View) {
	a.deleter.Execute(r, a.presenter.SetView(view))
}
func (a *Annotator) UpdateLabel(r updlbl.Request, view v.View) {
	a.labelUpdater.Execute(r, a.presenter.SetView(view))
}
func (a *Annotator) UpdateBox(r updbox.Request, view v.View) {
	a.boxUpdater.Execute(r, a.presenter.SetView(view))
}
func (a *Annotator) AddBox(r addbox.Request, view v.View) {
	a.boxAdder.Execute(r, a.presenter.SetView(view))
}
func (a *Annotator) AddLabel(r addlbl.Request, view v.View) {
	a.imageLabelAdder.Execute(r, a.presenter.SetView(view))
}

func (a *Annotator) Init(imageId string, collection string, view v.View) {
	scrollerState, err := a.scroller.Init(imageId, scr.WithCollection(collection))
	if err != nil {
		view.Error(err)
		return
	}
	view.SetScroller(p.MakeScrollerButtons(*scrollerState))

	a.presenter.SetView(view)
	a.imageReader.Execute(imread.Request{ImageId: imageId, Collection: collection},
		a.presenter)
	a.labelFetcher.Execute(a.presenter)
}

func NewAnnotator(scroller scr.Interface, imageMetaReader imread.Interface,
	boxAdder addbox.Interface, boxUpdater updbox.Interface, annotationDeleter del.Interface,
	labelFetcher fetchlbl.Interface, labelUpdater updlbl.Interface, imageLabelAdder addlbl.Interface,
	presenter p.Presenter) *Annotator {
	return &Annotator{
		scroller:        scroller,
		imageReader:     imageMetaReader,
		boxAdder:        boxAdder,
		boxUpdater:      boxUpdater,
		deleter:         annotationDeleter,
		labelFetcher:    labelFetcher,
		labelUpdater:    labelUpdater,
		imageLabelAdder: imageLabelAdder,
		presenter:       presenter,
	}
}
