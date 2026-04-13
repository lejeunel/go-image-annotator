package annotator

import (
	"github.com/lejeunel/go-image-annotator-v2/application/scroller"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	imread "github.com/lejeunel/go-image-annotator-v2/use-cases/image/read"
)

type Annotator struct {
	scroller          scroller.Interface
	imageReader       imread.Interface
	boxAdder          addbox.Interface
	boxUpdater        updbox.Interface
	annotationDeleter del.Interface
}

func (a *Annotator) DeleteAnnotation(r del.Request, p IAnnotatorPresenter) {
	a.annotationDeleter.Execute(r, p)
}

func (a *Annotator) UpdateBox(r updbox.Request, p IAnnotatorPresenter) {
	a.boxUpdater.Execute(r, p)
}

func (a *Annotator) AddBox(r addbox.Request, p IAnnotatorPresenter) {
	a.boxAdder.Execute(r, p)
}
func (a *Annotator) Start(imageId im.ImageId, collection string, p IAnnotatorPresenter) {
	scrollerState, err := a.scroller.Init(imageId, scroller.WithCollection(collection))
	if err != nil {
		p.Error(err)
		return
	}
	p.UpdateScroller(*scrollerState)
	a.imageReader.Execute(imread.Request{ImageId: imageId, Collection: collection}, p)

}

func NewAnnotator(scroller scroller.Interface, imageMetaReader imread.Interface,
	boxAdder addbox.Interface, boxUpdater updbox.Interface, annotationDeleter del.Interface) *Annotator {
	return &Annotator{
		scroller:          scroller,
		imageReader:       imageMetaReader,
		boxAdder:          boxAdder,
		boxUpdater:        boxUpdater,
		annotationDeleter: annotationDeleter}
}
