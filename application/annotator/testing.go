package annotator

import (
	scr "github.com/lejeunel/go-image-annotator-v2/application/scroller"
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

func (s *FakeScroller) Init(imageId im.ImageId, opts ...scr.Option) (*scr.ScrollerState, error) {
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

type FakePresenter struct {
	UpdatedScroller bool
	PresentedLabels bool
	PresentedImage  *im.Image
	AddedBox        *addbox.Response
	UpdatedBox      *updbox.Response
	DeletedBox      *del.Response
	GotErr          error
}

func (v *FakePresenter) Error(err error) {
	v.GotErr = err
}
func (v *FakePresenter) UpdateScroller(s scr.ScrollerState) {
	v.UpdatedScroller = true
}

func (v *FakePresenter) SuccessReadImage(i im.Image) {
	v.PresentedImage = &i
}
func (v *FakePresenter) SuccessAddBox(r addbox.Response) {
	v.AddedBox = &r
}
func (v *FakePresenter) SuccessUpdateBox(r updbox.Response) {
}
func (v *FakePresenter) SuccessDeleteAnnotation(r del.Response) {
}
func (v *FakePresenter) SuccessFetchLabels(r fetchlbl.Response) {
	v.PresentedLabels = true
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
	Got addbox.Request
}

func (b *FakeBoxAdder) Execute(r addbox.Request, o addbox.OutputPort) {
	b.Got = r
}

type FakeBoxUpdater struct {
	Got updbox.Request
}

func (b *FakeBoxUpdater) Execute(r updbox.Request, o updbox.OutputPort) {
	b.Got = r
}

type FakeAnnotationDeleter struct {
	Got del.Request
}

func (b *FakeAnnotationDeleter) Execute(r del.Request, o del.OutputPort) {
	b.Got = r
}

type FakeView struct {
	GotScrollerButtons ScrollerButtons
	GotErr             error
	GotImage           Image
	GotImageInfo       ImageInfo
}

func (s *FakeView) DrawScroller(buttons ScrollerButtons) {
	s.GotScrollerButtons = buttons
}
func (s *FakeView) Error(error)     {}
func (s *FakeView) DrawImage(Image) {}
func (s *FakeView) DrawImageInfo(i ImageInfo) {
	s.GotImageInfo = i
}
func (s *FakeView) AddBox(addbox.Response)        {}
func (s *FakeView) UpdateBox(updbox.Response)     {}
func (s *FakeView) DeleteAnnotation(del.Response) {}
func (s *FakeView) SetAvailableLabels([]string)   {}
