package annotator

import (
	scr "github.com/lejeunel/go-image-annotator-v2/application/scroller"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	fetchlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/label/fetch-all"
)

type IAnnotatorPresenter interface {
	UpdateScroller(scr.ScrollerState)
	SuccessReadImage(im.Image)
	SuccessAddBox(r addbox.Response)
	SuccessUpdateBox(r updbox.Response)
	SuccessDeleteAnnotation(r del.Response)
	SuccessFetchLabels(r fetchlbl.Response)
	Error(err error)
}

type AnnotatorPresenter struct {
	view AnnotatorView
}

func (p AnnotatorPresenter) UpdateScroller(s scr.ScrollerState) {
	buttons := ScrollerButtons{}
	if s.Next != nil {
		buttons.Next = ScrollerButton{IsActive: true,
			Text:       "Next",
			ImageId:    s.Next.ImageId.String(),
			Collection: s.Next.Collection}
	}
	if s.Previous != nil {
		buttons.Prev = ScrollerButton{IsActive: true,
			Text:       "Previous",
			ImageId:    s.Previous.ImageId.String(),
			Collection: s.Previous.Collection}
	}
	p.view.DrawScroller(buttons)
}
func (p AnnotatorPresenter) SuccessReadImage(im im.Image) {
	p.view.DrawImageInfo(NewImageInfo(im.Id, im.Collection.Name))
	p.view.DrawImage(NewImage(im.Id, im.Reader, im.Collection.Name, im.MIMEType))
}
func (p AnnotatorPresenter) SuccessFetchLabels(r fetchlbl.Response) {
	p.view.SetAvailableLabels(r.Labels)
}
func (p AnnotatorPresenter) SuccessAddBox(box addbox.Response) {
}
func (p AnnotatorPresenter) SuccessUpdateBox(box updbox.Response) {
}
func (p AnnotatorPresenter) SuccessDeleteAnnotation(a del.Response) {
}
func (p AnnotatorPresenter) Error(err error) {
}

func NewAnnotatorPresenter(view AnnotatorView) *AnnotatorPresenter {
	return &AnnotatorPresenter{view}
}
