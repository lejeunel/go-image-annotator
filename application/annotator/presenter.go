package annotator

import (
	scr "github.com/lejeunel/go-image-annotator-v2/application/scroller"
	a "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	fetchlbl "github.com/lejeunel/go-image-annotator-v2/use-cases/label/fetch-all"
)

type Presenter interface {
	SuccessReadImage(im.Image)
	SuccessAddBox(b a.BoundingBox)
	SuccessUpdateBox(r updbox.Response)
	SuccessDeleteAnnotation(r del.Response)
	SuccessFetchLabels(r fetchlbl.Response)
	Error(err error)
}

type AnnotatorPresenter struct {
	view View
}

func MakeScrollerButtons(s scr.ScrollerState) ScrollerButtons {
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
	return buttons
}

type StartPresenter struct {
	view View
}

func (p StartPresenter) SuccessReadImage(im im.Image) {
	p.view.DrawImageInfo(NewImageInfo(im.Id, im.Collection.Name))
	p.view.DrawImage(NewImage(im.Id, im.Reader, im.Collection.Name, im.MIMEType))
	p.view.DrawAnnotationList(im.BoundingBoxes)
}
func (p StartPresenter) SuccessFetchLabels(r fetchlbl.Response) {
	p.view.SetAvailableLabels(r.Labels)
}
func (p StartPresenter) Error(err error) {}

type AddBoxPresenter struct {
	view View
}

func (p AddBoxPresenter) SuccessAddBox(b a.BoundingBox) {
	p.view.AddBox(b)
}
func (p AddBoxPresenter) Error(err error) {}

type RemoveAnnotationPresenter struct {
	view View
}

func (p RemoveAnnotationPresenter) SuccessDeleteAnnotation(r del.Response) {
	p.view.DeleteAnnotation(r)
}
func (p RemoveAnnotationPresenter) Error(err error) {}
