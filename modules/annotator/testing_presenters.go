package annotator

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	scr "github.com/lejeunel/go-image-annotator/modules/scroller"
	addbox "github.com/lejeunel/go-image-annotator/use-cases/annotate/add-bbox"
	rmlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/remove"
	updlbl "github.com/lejeunel/go-image-annotator/use-cases/annotate/update-label"
	fetchlbl "github.com/lejeunel/go-image-annotator/use-cases/label/fetch-all"
)

type FakeScrollerPresenter struct {
}

func (p *FakeScrollerPresenter) SuccessInitScroller(scr.ScrollerState) {
}
func (p FakeScrollerPresenter) Error(error) {}

type FakeLabelFetchPresenter struct {
	Called bool
}

func (p *FakeLabelFetchPresenter) SuccessFetchLabels(fetchlbl.Response) {
	p.Called = true
}
func (p FakeLabelFetchPresenter) Error(error) {}

type FakeImageReadPresenter struct {
	Called bool
}

func (p *FakeImageReadPresenter) SuccessReadImage(im.Image) {
	p.Called = true
}
func (p FakeImageReadPresenter) Error(error) {}

type FakeAddBoxPresenter struct {
	Called bool
}

func (p *FakeAddBoxPresenter) SuccessAddBox(addbox.Response) {
	p.Called = true
}
func (p FakeAddBoxPresenter) Error(error) {}

type FakeUpdateLabelPresenter struct {
	Called bool
}

func (p *FakeUpdateLabelPresenter) SuccessUpdateLabel(updlbl.Response) {
	p.Called = true
}
func (p FakeUpdateLabelPresenter) Error(error) {}

type FakeRemoveLabelPresenter struct {
	Called bool
}

func (p *FakeRemoveLabelPresenter) SuccessDeleteAnnotation(rmlbl.Response) {
	p.Called = true
}
func (p FakeRemoveLabelPresenter) Error(error) {}
