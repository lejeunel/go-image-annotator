package annotator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"embed"

	an "github.com/lejeunel/go-image-annotator/adapters/web/annotator/annotorious"
	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	v "github.com/lejeunel/go-image-annotator/modules/annotator/view"
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
	image                *v.Image
	boxes                []v.BoundingBox
	imageLabels          []v.ImageLabel
	imageInfo            *v.ImageInfo
	availableLabels      []string
	availableImageLabels []string
	scrollerButtons      v.ScrollerButtons
	err                  error
	PageBuilder          b.PageBuilder
}

func (v *AnnotationView) SetScroller(buttons v.ScrollerButtons) {
	v.scrollerButtons = buttons
}
func (v *AnnotationView) SetAnnotations(boxes []v.BoundingBox, imageLabels []v.ImageLabel) {
	v.boxes = boxes
	v.imageLabels = imageLabels
}
func (v *AnnotationView) SetAvailableLabels(labels []string) {
	v.availableLabels = labels
}
func (v *AnnotationView) SetAvailableImageLabels(labels []string) {
	v.availableImageLabels = labels
}
func (v *AnnotationView) SetImageInfo(info v.ImageInfo) {
	v.imageInfo = &info
}
func (v *AnnotationView) SetImage(image v.Image) {
	v.image = &image
}
func (v *AnnotationView) AddBox(b v.BoundingBox)     {}
func (v *AnnotationView) AddLabel(l v.ImageLabel)    {}
func (v *AnnotationView) UpdateBox(b v.BoundingBox)  {}
func (v *AnnotationView) UpdateLabel(a v.Annotation) {}
func (v *AnnotationView) DeleteAnnotation(string)    {}
func (v *AnnotationView) Error(err error) {
	v.err = err
}
func (v *AnnotationView) RenderAnnotationList(w http.ResponseWriter) {
	v.AnnotationsListView.Build(v.boxes, v.imageLabels, v.availableLabels).Render(w)
}
func (v *AnnotationView) RenderAll(w http.ResponseWriter) {

	if v.err != nil {
		http.Error(w, v.err.Error(), http.StatusBadRequest)
		return
	}

	if v.image != nil {
		v.render(w)
	}

}
func (v *AnnotationView) RenderAnnotations(w http.ResponseWriter) {
	if v.err != nil {
		http.Error(w, v.err.Error(), http.StatusBadRequest)
		return
	}
	boxes := an.ConvertBoxesToAnnotorious(v.boxes)
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

func (v *AnnotationView) ShapeSelector() Node {
	classInactive := `flex items-center whitespace-nowrap rounded-radius bg-surface-alt border border-surface-alt px-2 py-1 text-sm font-medium tracking-wide text-on-surface-strong transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-surface-alt active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed dark:bg-surface-dark-alt dark:border-surface-dark-alt dark:text-on-surface-dark-strong dark:focus-visible:outline-surface-dark-alt`
	classActive := `flex items-center whitespace-nowrap rounded-radius bg-primary border border-surface-primary px-2 py-1 text-sm font-medium tracking-wide text-on-primary transition hover:opacity-75 text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary active:opacity-100 active:outline-offset-0 disabled:opacity-75 disabled:cursor-not-allowed dark:bg-primary-dark dark:border-primary-dark dark:text-on-primary-dark dark:focus-visible:outline-primary-dark`

	return Div(
		Attr("x-data", "{ active: 'rectangle'}"),
		Class("flex gap-2 pb-2"),
		Button(
			Attr("x-bind:class", fmt.Sprintf(`{'%v': active === 'rectangle', '%v': active !== 'rectangle'}`, classActive, classInactive)),
			Attr("@click", "AnnotatorModule.drawRectangle(); active = 'rectangle';"),
			Raw(RectangleIcon), Div(Class("ml-1"), Text("Rectangle"))),
		Button(
			Attr("x-bind:class", fmt.Sprintf(`{'%v': active === 'polygon', '%v': active !== 'polygon'}`, classActive, classInactive)),
			Attr("@click", "AnnotatorModule.drawPolygon(); active = 'polygon';"),
			Raw(PolygonIcon), Div(Class("ml-1"), Text("Polygon"))))
}

func AnnotoriousLib() []Node {
	var scripts []Node
	scripts = append(scripts, Script(Defer(), Src("/static/annotorious.js")))
	scripts = append(scripts, Link(Href("/static/annotorious.css"), Rel("stylesheet")))
	return scripts
}

func (v *AnnotationView) render(w http.ResponseWriter) {
	pb := v.PageBuilder

	script, err := MakeAnnotoriousScript(v.image.Id, v.image.Collection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	regionLabelModal, _ := makeLabelModal(v.availableLabels, RegionLabelModal)
	imageLabelModal, _ := makeLabelModal(v.availableLabels, ImageLabelModal)
	pb.AddScripts(AnnotoriousLib()...)
	pb.AddScripts(*script)
	pb.AddScripts(Raw(*regionLabelModal))
	pb.AddScripts(Raw(*imageLabelModal))

	pb.SetContent(
		Table(
			Tr(Div(Class("flex"), v.ScrollerView.Render(v.scrollerButtons), v.ShapeSelector())),
			Tr(Td(Table(
				Tr(Td(Class("align-top"), v.ImageView.Build(*v.image)),
					Td(Class("align-top pl-2"),
						Div(Class("pb-2"), v.ImageInfosView.Build(*v.imageInfo)),
						Div(ID("annotation-list"), v.AnnotationsListView.Build(v.boxes, v.imageLabels, v.availableLabels)))),
			),
			))))
	pb.Render(w)
}

func NewAnnotationView(pageBuilder b.PageBuilder) *AnnotationView {
	return &AnnotationView{
		ImageView:      ImageView{},
		ImageInfosView: ImageInfosView{},
		ScrollerView:   ScrollerView{},
		PageBuilder:    *pageBuilder.SetTitle("Image"),
	}
}
