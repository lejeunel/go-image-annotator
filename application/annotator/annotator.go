package annotator

import (
	"github.com/lejeunel/go-image-annotator-v2/application/scroller"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	imread "github.com/lejeunel/go-image-annotator-v2/use-cases/image/read"
	fetchlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/label/fetch-all"
)

type Annotator struct {
	scroller     scroller.Interface
	imageReader  imread.Interface
	boxAdder     addbox.Interface
	boxUpdater   updbox.Interface
	deleter      del.Interface
	labelFetcher fetchlbl.Interface
}

func (a *Annotator) DeleteAnnotation(r del.Request, view View) {
	a.deleter.Execute(r, RemoveAnnotationPresenter{view})
}
func (a *Annotator) UpdateBox(r updbox.Request) {
	// a.boxUpdater.Execute(r, p)
}
func (a *Annotator) AddBox(r addbox.Request, view View) {
	a.boxAdder.Execute(r, AddBoxPresenter{view})
}
func (a *Annotator) Init(imageId im.ImageId, collection string, view View) {
	scrollerState, err := a.scroller.Init(imageId, scroller.WithCollection(collection))
	if err != nil {
		view.Error(err)
		return
	}
	view.DrawScroller(MakeScrollerButtons(*scrollerState))

	presenter := StartPresenter{view}
	a.imageReader.Execute(imread.Request{ImageId: imageId, Collection: collection},
		presenter)
	a.labelFetcher.Execute(presenter)
}

func NewAnnotator(scroller scroller.Interface, imageMetaReader imread.Interface,
	boxAdder addbox.Interface, boxUpdater updbox.Interface, annotationDeleter del.Interface,
	labelFetcher fetchlbl.Interface) *Annotator {
	return &Annotator{
		scroller:     scroller,
		imageReader:  imageMetaReader,
		boxAdder:     boxAdder,
		boxUpdater:   boxUpdater,
		deleter:      annotationDeleter,
		labelFetcher: labelFetcher,
	}
}
