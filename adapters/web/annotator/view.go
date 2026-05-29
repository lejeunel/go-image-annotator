package annotator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"embed"

	v "github.com/lejeunel/go-image-annotator/app/annotator/view"
	html "github.com/lejeunel/go-image-annotator/shared/html"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

//go:embed templates/*
var templatesFiles embed.FS

type AnnotationView struct {
	ImageView
	ImageInfosView
	AnnotationsListView
	ScrollerView
	image           *v.Image
	boxes           []*v.BoundingBox
	imageLabels     []*v.ImageLabel
	imageInfo       *v.ImageInfo
	availableLabels []string
	scrollerButtons v.ScrollerButtons
	doDrawImage     bool
	addedBox        *v.BoundingBox
	addedLabel      *v.ImageLabel
	err             error
}

func (v *AnnotationView) DrawScroller(buttons v.ScrollerButtons) {
	v.scrollerButtons = buttons
}
func (v *AnnotationView) DrawImage(image v.Image) {
	v.image = &image
	v.doDrawImage = true
}
func (v *AnnotationView) DrawImageInfo(info v.ImageInfo) {
	v.imageInfo = &info
}
func (v *AnnotationView) DrawAnnotationList(boxes []*v.BoundingBox, imageLabels []*v.ImageLabel) {
	v.boxes = boxes
	v.imageLabels = imageLabels
}
func (v *AnnotationView) SetAvailableLabels(labels []string) {
	v.availableLabels = labels

}
func (v *AnnotationView) AddBox(b v.BoundingBox) {
	v.addedBox = &b
}
func (v *AnnotationView) AddLabel(l v.ImageLabel) {
	v.addedLabel = &l
}
func (v *AnnotationView) UpdateBox(b v.BoundingBox) {
}
func (v *AnnotationView) UpdateLabel(a v.Annotation) {
}
func (v *AnnotationView) DeleteAnnotation(string) {
}
func (v *AnnotationView) Error(err error) {
	v.err = err
}
func (v *AnnotationView) RenderAnnotationList(w http.ResponseWriter) {
	v.AnnotationsListView.Build(v.boxes, v.imageLabels).Render(w)
}
func (v *AnnotationView) RenderAll(w http.ResponseWriter) {

	if v.err != nil {
		http.Error(w, v.err.Error(), http.StatusBadRequest)
		return
	}

	if v.doDrawImage {
		v.renderImage(w)
	}

}
func (v *AnnotationView) RenderAnnotations(w http.ResponseWriter) {
	if v.err != nil {
		http.Error(w, v.err.Error(), http.StatusBadRequest)
		return
	}
	boxes := ConvertToAnnotorious(v.boxes)
	data, err := json.Marshal(boxes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func (v *AnnotationView) renderImage(w http.ResponseWriter) {
	b := html.NewTitledPageBuilder("Image")
	script, err := MakeAnnotoriousScript(v.image.Id, v.image.Collection)
	if err != nil {
		b.SetError(err).Render(w)
		return
	}
	labelModal, err := makeLabelModal(v.availableLabels)
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
				Tr(Td(Class("align-top"), v.ImageView.Build(*v.image)),
					Td(Class("align-top pl-2"),
						Div(Class("pb-2"), v.ImageInfosView.Build(*v.imageInfo)),
						Div(ID("annotation-list"), v.AnnotationsListView.Build(v.boxes, v.imageLabels)))),
			),
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
