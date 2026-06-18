package annotator

import (
	"context"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-polygon"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	updpoly "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-polygon"
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
	polygonUpdater  updpoly.Interface
	imageLabelAdder addlbl.Interface
	boxUpdater      updbox.Interface
	deleter         del.Interface
	labelFetcher    fetchlbl.Interface
	labelUpdater    updlbl.Interface
}

func (a *Annotator) Init(ctx context.Context, imageId string, collection string,
	oim imread.OutputPort, olbl fetchlbl.OutputPort, oscr scr.OutputPort) {
	a.scroller.Init(imageId, oscr,
		scr.WithCollection(collection),
		scr.WithOrdering(im.OrderingParams{IngestTime: true}))
	a.ReadImage(imageId, collection, oim)
	a.labelFetcher.Execute(ctx, olbl)
}
func (a *Annotator) ReadImage(imageId string, collection string, o imread.OutputPort) {
	a.imageReader.Execute(imread.Request{ImageId: imageId, Collection: collection}, o)
}

func (a *Annotator) DeleteAnnotation(ctx context.Context, r del.Request, o del.OutputPort) {
	a.deleter.Execute(ctx, r, o)
}
func (a *Annotator) UpdateLabel(ctx context.Context, r updlbl.Request, o updlbl.OutputPort) {
	a.labelUpdater.Execute(ctx, r, o)
}
func (a *Annotator) UpdateBox(ctx context.Context, r updbox.Request, o updbox.OutputPort) {
	a.boxUpdater.Execute(ctx, r, o)
}
func (a *Annotator) UpdatePolygon(ctx context.Context, r updpoly.Request, o updpoly.OutputPort) {
	a.polygonUpdater.Execute(ctx, r, o)
}
func (a *Annotator) AddBox(ctx context.Context, r addbox.Request, o addbox.OutputPort) {
	a.boxAdder.Execute(ctx, r, o)
}
func (a *Annotator) AddPolygon(ctx context.Context, r addpoly.Request, o addpoly.OutputPort) {
	a.polygonAdder.Execute(ctx, r, o)
}
func (a *Annotator) AddLabel(ctx context.Context, r addlbl.Request, o addlbl.OutputPort) {
	a.imageLabelAdder.Execute(ctx, r, o)
}

func NewAnnotator(scroller scr.Interface, imageMetaReader imread.Interface,
	boxAdder addbox.Interface, boxUpdater updbox.Interface, polygonAdder addpoly.Interface,
	polygonUpdater updpoly.Interface, annotationDeleter del.Interface,
	labelFetcher fetchlbl.Interface, labelUpdater updlbl.Interface, imageLabelAdder addlbl.Interface) Annotator {
	return Annotator{
		scroller:        scroller,
		imageReader:     imageMetaReader,
		boxAdder:        boxAdder,
		boxUpdater:      boxUpdater,
		polygonAdder:    polygonAdder,
		polygonUpdater:  polygonUpdater,
		deleter:         annotationDeleter,
		labelFetcher:    labelFetcher,
		labelUpdater:    labelUpdater,
		imageLabelAdder: imageLabelAdder,
	}
}
