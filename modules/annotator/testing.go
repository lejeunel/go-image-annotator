package annotator

import (
	"context"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	addlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/assign-label"
	updbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
	imread "github.com/lejeunel/go-image-annotator/use-cases/image/read"
	fetchlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
)

type FakeScroller struct {
	Err       error
	ErrOnInit bool
	IsInit    bool
}

func (s *FakeScroller) Init(imageId string, opts ...scr.Option) (*scr.ScrollerState, error) {
	if s.ErrOnInit {
		return nil, s.Err
	}
	s.IsInit = true
	return &scr.ScrollerState{}, nil
}

type FakeLabelFetcher struct{}

func (f *FakeLabelFetcher) Execute(ctx context.Context, o fetchlbl.OutputPort) {
	o.SuccessFetchLabels(fetchlbl.Response{Labels: []string{"a-label"}})
}

type FakeImageReader struct {
	Return *im.Image
}

func (b *FakeImageReader) Execute(r imread.Request, o imread.OutputPort) {
	o.SuccessReadImage(*b.Return)
}

type FakeLabelAdder struct {
	Returns addlbl.Response
}

func (b *FakeLabelAdder) Execute(ctx context.Context, r addlbl.Request, o addlbl.OutputPort) {
	o.SuccessAddLabel(b.Returns)
}

type FakeBoxAdder struct {
	Returns addbox.Response
}

func (b *FakeBoxAdder) Execute(c context.Context, r addbox.Request, o addbox.OutputPort) {
	o.SuccessAddBox(b.Returns)
}

type FakeBoxUpdater struct {
	Returns *updbox.Response
}

func (b *FakeBoxUpdater) Execute(c context.Context, r updbox.Request, o updbox.OutputPort) {
	o.SuccessUpdateBox(*b.Returns)
}

type FakeLabelUpdater struct {
}

func (b *FakeLabelUpdater) Execute(ctx context.Context, r updlbl.Request, o updlbl.OutputPort) {
	o.SuccessUpdateLabel(updlbl.Response{})
}

type FakeAnnotationDeleter struct {
	Returns del.Response
}

func (b *FakeAnnotationDeleter) Execute(c context.Context, r del.Request, o del.OutputPort) {
	o.SuccessDeleteAnnotation(b.Returns)
}

type FakeView struct {
	GotScrollerButtons       *v.ScrollerButtons
	GotErr                   error
	GotImage                 *v.Image
	GotImageInfo             *v.ImageInfo
	GotAvailableRegionLabels *[]string
	GotAvailableImageLabels  *[]string
	GotAnnotationIds         *[]string
	RemovedAnnotationId      *string
	UpdatedBoxId             *string
	UpdatedAnnotation        *v.Annotation
}

func (s *FakeView) SetScroller(buttons v.ScrollerButtons) {
	s.GotScrollerButtons = &buttons
}
func (s *FakeView) Error(error) {}
func (s *FakeView) SetImage(i v.Image) {
	s.GotImage = &i
}
func (s *FakeView) SetImageInfo(i v.ImageInfo) {
	s.GotImageInfo = &i
}
func (s *FakeView) SetAnnotations(boxes []v.BoundingBox, labels []v.ImageLabel) {
	ids := []string{}
	for _, b := range boxes {
		ids = append(ids, b.Id)
	}
	for _, l := range labels {
		ids = append(ids, l.Id)
	}
	s.GotAnnotationIds = &ids
}

func (s *FakeView) UpdateBox(b v.BoundingBox) {
	s.UpdatedBoxId = &b.Id
}
func (s *FakeView) UpdateLabel(a v.Annotation) {
	s.UpdatedAnnotation = &a
}

func (s *FakeView) DeleteAnnotation(id string) {
	s.RemovedAnnotationId = &id
}
func (s *FakeView) SetAvailableLabels(l []string) {
	s.GotAvailableRegionLabels = &l
}

func (s *FakeView) SetAvailableImageLabels(l []string) {
	s.GotAvailableImageLabels = &l
}
