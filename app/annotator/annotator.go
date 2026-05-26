package annotator

import (
	p "github.com/lejeunel/go-image-annotator-v2/app/annotator/presenters"
	scr "github.com/lejeunel/go-image-annotator-v2/app/annotator/scroller"
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/update-label"
	imread "github.com/lejeunel/go-image-annotator-v2/use-cases/image/read"
	fetchlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/label/fetch-all"
)

type Annotator struct {
	scroller     scr.Interface
	imageReader  imread.Interface
	boxAdder     addbox.Interface
	boxUpdater   updbox.Interface
	deleter      del.Interface
	labelFetcher fetchlbl.Interface
	labelUpdater updlbl.Interface
}

func (a *Annotator) DeleteAnnotation(r del.Request, view v.View) {
	a.deleter.Execute(r, p.RemoveAnnotationPresenter{View: view})
}
func (a *Annotator) UpdateLabelOfAnnotation(r updlbl.Request, view v.View) {
	a.labelUpdater.Execute(r, p.UpdateLabelOfAnnotationPresenter{View: view})
}
func (a *Annotator) UpdateBox(r updbox.Request, view v.View) {
	a.boxUpdater.Execute(r, p.UpdateBoxPresenter{View: view})
}
func (a *Annotator) AddBox(r addbox.Request, view v.View) {
	a.boxAdder.Execute(r, p.AddBoxPresenter{View: view})
}
func (a *Annotator) Init(imageId string, collection string, view v.View) {
	scrollerState, err := a.scroller.Init(imageId, scr.WithCollection(collection))
	if err != nil {
		view.Error(err)
		return
	}
	view.DrawScroller(p.MakeScrollerButtons(*scrollerState))

	presenter := p.InitPresenter{View: view}
	a.imageReader.Execute(imread.Request{ImageId: imageId, Collection: collection},
		presenter)
	a.labelFetcher.Execute(presenter)
}

func NewAnnotator(scroller scr.Interface, imageMetaReader imread.Interface,
	boxAdder addbox.Interface, boxUpdater updbox.Interface, annotationDeleter del.Interface,
	labelFetcher fetchlbl.Interface, labelUpdater updlbl.Interface) *Annotator {
	return &Annotator{
		scroller:     scroller,
		imageReader:  imageMetaReader,
		boxAdder:     boxAdder,
		boxUpdater:   boxUpdater,
		deleter:      annotationDeleter,
		labelFetcher: labelFetcher,
		labelUpdater: labelUpdater,
	}
}
