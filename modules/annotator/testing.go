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
)

type FakeScroller struct {
	Err       error
	ErrOnInit bool
	IsInit    bool
}

func (s *FakeScroller) Init(imageId string, o scr.OutputPort, opts ...scr.Option) {
	if s.ErrOnInit {
		return
	}
	s.IsInit = true
}

type FakeLabelFetcher struct{}

func (f *FakeLabelFetcher) Execute(ctx context.Context, o fetchlbl.OutputPort) {
	o.SuccessFetchLabels(fetchlbl.Response{Labels: []string{"a-label"}})
}

type FakeImageReader struct {
}

func (b *FakeImageReader) Execute(r imread.Request, o imread.OutputPort) {
	o.SuccessReadImage(im.Image{})
}

type FakeLabelAdder struct {
}

func (b *FakeLabelAdder) Execute(ctx context.Context, r addlbl.Request, o addlbl.OutputPort) {
	o.SuccessAddLabel(addlbl.Response{})
}

type FakePolygonAdder struct {
	Returns addpoly.Response
}

func (b *FakePolygonAdder) Execute(c context.Context, r addpoly.Request, o addpoly.OutputPort) {
	o.SuccessAddPolygon(addpoly.Response{})
}

type FakePolygonUpdater struct {
}

func (b *FakePolygonUpdater) Execute(c context.Context, r updpoly.Request, o updpoly.OutputPort) {
	o.SuccessUpdatePolygon(updpoly.Response{})
}

type FakeBoxAdder struct {
}

func (b *FakeBoxAdder) Execute(c context.Context, r addbox.Request, o addbox.OutputPort) {
	o.SuccessAddBox(addbox.Response{})
}

type FakeBoxUpdater struct {
}

func (b *FakeBoxUpdater) Execute(c context.Context, r updbox.Request, o updbox.OutputPort) {
	o.SuccessUpdateBox(updbox.Response{})
}

type FakeLabelUpdater struct {
}

func (b *FakeLabelUpdater) Execute(ctx context.Context, r updlbl.Request, o updlbl.OutputPort) {
	o.SuccessUpdateLabel(updlbl.Response{})
}

type FakeAnnotationDeleter struct {
}

func (b *FakeAnnotationDeleter) Execute(c context.Context, r del.Request, o del.OutputPort) {
	o.SuccessDeleteAnnotation(del.Response{})
}
