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
	imread "github.com/lejeunel/go-image-annotator/use-cases/image/find"
	fetchlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
	addmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/add"
	delmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/delete"
	listmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/list"
	readmd "github.com/lejeunel/go-image-annotator/use-cases/metadata/read"
)

type Annotator struct {
	scroller         scr.Interface
	readImage        imread.Interface
	AddBox           addbox.Interface
	AddPolygon       addpoly.Interface
	UpdatePolygon    updpoly.Interface
	AddLabel         addlbl.Interface
	UpdateBox        updbox.Interface
	DeleteAnnotation del.Interface
	FetchLabels      fetchlbl.Interface
	UpdateLabel      updlbl.Interface
	AddMetaData      addmd.Interface
	ListMetaData     listmd.Interface
	ReadMetaData     readmd.Interface
	DeleteMetaData   delmd.Interface
}

func (a *Annotator) Init(ctx context.Context, imageId string, collection string,
	oim imread.OutputPort, olbl fetchlbl.OutputPort, oscr scr.OutputPort,
) {
	a.scroller.Init(imageId, oscr,
		scr.WithCollection(collection),
		scr.WithOrdering(im.Ordering{IngestTime: true}))
	a.ReadImage(imageId, collection, oim)
	a.FetchLabels.Execute(ctx, olbl)
}

func (a *Annotator) ReadImage(imageId string, collection string, o imread.OutputPort) {
	a.readImage.Execute(imread.Request{ImageId: imageId, Collection: collection}, o)
}

func NewAnnotator(
	scroller scr.Interface,
	imageMetaReader imread.Interface,
	boxAdder addbox.Interface,
	boxUpdater updbox.Interface,
	polygonAdder addpoly.Interface,
	polygonUpdater updpoly.Interface,
	annotationDeleter del.Interface,
	labelFetcher fetchlbl.Interface,
	labelUpdater updlbl.Interface,
	imageLabelAdder addlbl.Interface,
	metaAdder addmd.Interface,
	metaList listmd.Interface,
	metaRead readmd.Interface,
	metaDelete delmd.Interface,
) Annotator {
	return Annotator{
		scroller:         scroller,
		readImage:        imageMetaReader,
		AddBox:           boxAdder,
		UpdateBox:        boxUpdater,
		AddPolygon:       polygonAdder,
		UpdatePolygon:    polygonUpdater,
		DeleteAnnotation: annotationDeleter,
		FetchLabels:      labelFetcher,
		UpdateLabel:      labelUpdater,
		AddLabel:         imageLabelAdder,
		AddMetaData:      metaAdder,
		ListMetaData:     metaList,
		ReadMetaData:     metaRead,
		DeleteMetaData:   metaDelete,
	}
}
