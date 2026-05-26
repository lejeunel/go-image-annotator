package annotator

import (
	scr "github.com/lejeunel/go-image-annotator-v2/app/annotator/scroller"
	v "github.com/lejeunel/go-image-annotator-v2/app/annotator/view"
	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	imread "github.com/lejeunel/go-image-annotator-v2/use-cases/image/read"
	fetchlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/label/fetch-all"
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

func (f *FakeLabelFetcher) Execute(o fetchlbl.OutputPort) {
	o.SuccessFetchLabels(fetchlbl.Response{Labels: []string{"a-label"}})
}

type FakeImageReader struct {
	Got    imread.Request
	Return *im.Image
}

func (b *FakeImageReader) Execute(r imread.Request, o imread.OutputPort) {
	o.SuccessReadImage(*b.Return)
	b.Got = r
}

type FakeBoxAdder struct {
	Got     addbox.Request
	Returns an.BoundingBox
}

func (b *FakeBoxAdder) Execute(r addbox.Request, o addbox.OutputPort) {
	b.Got = r
	o.SuccessAddBox(b.Returns)
}

type FakeBoxUpdater struct {
	Got     updbox.Request
	Returns *updbox.Response
}

func (b *FakeBoxUpdater) Execute(r updbox.Request, o updbox.OutputPort) {
	b.Got = r
	o.SuccessUpdateBox(*b.Returns)
}

type FakeAnnotationDeleter struct {
	Got     del.Request
	Returns del.Response
}

func (b *FakeAnnotationDeleter) Execute(r del.Request, o del.OutputPort) {
	b.Got = r
	o.SuccessDeleteAnnotation(b.Returns)
}

type FakeView struct {
	GotScrollerButtons         *v.ScrollerButtons
	GotErr                     error
	GotBox                     *v.BoundingBox
	GotImage                   *v.Image
	GotImageInfo               *v.ImageInfo
	GotLabels                  *[]string
	RemovedAnnotationId        *an.AnnotationId
	UpdatedBoxId               *an.AnnotationId
	UpdatedLabelOfAnnotationId *an.AnnotationId
	UpdatedLabelOfAnnotation   string
}

func (s *FakeView) DrawScroller(buttons v.ScrollerButtons) {
	s.GotScrollerButtons = &buttons
}
func (s *FakeView) Error(error) {}
func (s *FakeView) DrawImage(i v.Image) {
	s.GotImage = &i
}
func (s *FakeView) DrawImageInfo(i v.ImageInfo) {
	s.GotImageInfo = &i
}
func (s *FakeView) DrawAnnotationList([]*v.BoundingBox, []*v.ImageLabel) {}
func (s *FakeView) AddBox(box v.BoundingBox) {
	s.GotBox = &box
}
func (s *FakeView) UpdateBox(r updbox.Response) {
	s.UpdatedBoxId = &r.AnnotationId
}
func (s *FakeView) DeleteAnnotation(r del.Response) {
	s.RemovedAnnotationId = &r.Id
}
func (s *FakeView) SetAvailableLabels(l []string) {
	s.GotLabels = &l
}
