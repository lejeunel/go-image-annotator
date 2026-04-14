package annotator

import (
	"bytes"
	"fmt"
	"io"

	"embed"

	"text/template"

	a "github.com/lejeunel/go-image-annotator-v2/application/annotator"
	an "github.com/lejeunel/go-image-annotator-v2/entities/annotation"
	html "github.com/lejeunel/go-image-annotator-v2/shared/html"
	updbox "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/modify-bbox"
	del "github.com/lejeunel/go-image-annotator-v2/use-cases/annotate/remove"
	. "maragu.dev/gomponents"
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
	availableLabels []string
	scrollerButtons a.ScrollerButtons
	doDrawImage     bool
	doAddBox        bool
	addedBox        an.BoundingBox
	err             error
}

func (v *AnnotationView) DrawScroller(buttons a.ScrollerButtons) {
	v.scrollerButtons = buttons
}

func (v *AnnotationView) DrawImage(image a.Image) {
	v.image = &image
	v.doDrawImage = true
}

func (v *AnnotationView) DrawImageInfo(info a.ImageInfo) {
	v.imageInfo = &info
}
func (v *AnnotationView) SetAvailableLabels(labels []string) {
	v.availableLabels = labels

}
func (v *AnnotationView) AddBox(b an.BoundingBox) {
	v.doAddBox = true
	v.addedBox = b
}
func (v *AnnotationView) UpdateBox(r updbox.Response) {
}
func (v *AnnotationView) DeleteAnnotation(r del.Response) {
}

func (v *AnnotationView) Error(err error) {
	v.err = err
}
func (v *AnnotationView) makeLabelModal(labels []string) (string, error) {
	tLabelModal, err := template.New("labelModal").ParseFS(templatesFiles,
		"templates/label_selector.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := struct {
		Labels []string
	}{labels}
	err = tLabelModal.ExecuteTemplate(&buf, "labelModal",
		data)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (v *AnnotationView) Render(w io.Writer) {

	if v.err != nil {
		html.NewPageBuilder().SetError(v.err).Render(w)
	}

	if v.doDrawImage {
		v.renderImage(w)
	}

	if v.doAddBox {
		fmt.Println("in view: added box: ", v.addedBox)
	}

}

func (v *AnnotationView) renderImage(w io.Writer) {
	b := html.NewTitledPageBuilder("Image")
	script, err := MakeAnnotoriousScript(v.image.Id, v.image.Collection)
	if err != nil {
		b.SetError(err).Render(w)
		return
	}
	labelModal, err := v.makeLabelModal(v.availableLabels)
	if err != nil {
		b.SetError(fmt.Errorf("building label model: %w", err)).Render(w)
		return
	}
	b.AddScripts(html.AnnotoriousLib()...)
	b.AddScripts(*script)
	b.AddScripts(Raw(labelModal))

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
