package annotator

import (
	"io"

	"embed"

	a "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	html "github.com/lejeunel/go-image-annotator-v2/shared/html"
	addbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/add-bbox"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	. "maragu.dev/gomponents/html"
)

//go:embed templates/*
var templatesFiles embed.FS

type AnnotationView struct {
	ImageView       ImageView
	ImageInfosView  ImageInfosView
	ScrollerView    ScrollerView
	image           *a.Image
	imageInfo       *a.ImageInfo
	scrollerButtons a.ScrollerButtons
	err             error
}

func (v *AnnotationView) DrawScroller(buttons a.ScrollerButtons) {
	v.scrollerButtons = buttons
}

func (v *AnnotationView) DrawImage(image a.Image) {
	v.image = &image
}

func (v *AnnotationView) DrawImageInfo(info a.ImageInfo) {
	v.imageInfo = &info
}

func (v *AnnotationView) AddBox(r addbox.Response) {
}
func (v *AnnotationView) UpdateBox(r updbox.Response) {
}
func (v *AnnotationView) DeleteAnnotation(r del.Response) {
}

func (v *AnnotationView) Error(err error) {
	v.err = err
}

func (v *AnnotationView) Render(w io.Writer) {

	if v.err != nil {
		html.NewPageBuilder().SetError(v.err).Render(w)
	}

	b := html.NewTitledPageBuilder("Image")
	script, err := MakeAnnotoriousScript(v.image.Id, v.image.Collection)
	if err != nil {
		b.SetError(err).Render(w)
		return
	}
	b.AddScripts(html.AnnotoriousLib()...)
	b.AddScripts(*script)

	b.SetContent(
		Table(
			Tr(Td(v.ScrollerView.Render(v.scrollerButtons))),
			Tr(Td(Table(
				Tr(Td(v.ImageView.Render(*v.image)),
					Td(Class("align-top pl-2"), v.ImageInfosView.Render(*v.imageInfo)))),
			))))
	b.Render(w)

}

func NewAnnotationView() *AnnotationView {
	return &AnnotationView{
		ImageView:      ImageView{},
		ImageInfosView: ImageInfosView{},
		ScrollerView:   ScrollerView{},
	}
}
