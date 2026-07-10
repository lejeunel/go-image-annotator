package annotator

import (
	"fmt"
	"net/http"

	"embed"

	b "github.com/lejeunel/go-image-annotator/adapters/web/builders"
	ic "github.com/lejeunel/go-image-annotator/adapters/web/icons"
	s "github.com/lejeunel/go-image-annotator/adapters/web/styles"
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
	polygons             []v.Polygon
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
func (v *AnnotationView) SetAnnotations(boxes []v.BoundingBox, polygons []v.Polygon, imageLabels []v.ImageLabel) {
	v.boxes = boxes
	v.polygons = polygons
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
func (v *AnnotationView) Error(err error) {
	v.err = err
}
func (v *AnnotationView) RenderAnnotationList(w http.ResponseWriter) {
	v.AnnotationsListView.Build(v.boxes, v.polygons, v.imageLabels, v.availableLabels).Render(w)
}
func (v *AnnotationView) Render(w http.ResponseWriter) {

	if v.err != nil {
		http.Error(w, v.err.Error(), http.StatusBadRequest)
		return
	}

	if v.image != nil {
		v.render(w)
	}
}

func (v *AnnotationView) ShapeSelector() Node {
	return Div(
		Attr("x-data", "{ active: 'rectangle'}"),
		Class("flex gap-2 pb-2"),
		Button(
			Attr("x-bind:class", fmt.Sprintf(`{'%v': active === 'rectangle', '%v': active !== 'rectangle'}`, s.PrimaryButton, s.InactiveButton)),
			Attr("@click", "AnnotatorModule.drawRectangle(); active = 'rectangle';"),
			Raw(ic.BoundingBoxIcon), Div(Class("ml-1"), Text("Rectangle"))),
		Button(
			Attr("x-bind:class", fmt.Sprintf(`{'%v': active === 'polygon', '%v': active !== 'polygon'}`, s.PrimaryButton, s.InactiveButton)),
			Attr("@click", "AnnotatorModule.drawPolygon(); active = 'polygon';"),
			Raw(ic.PolygonIcon), Div(Class("ml-1"), Text("Polygon"))))
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
	pb.AddScripts(AnnotoriousLib()...)
	pb.AddScripts(*script)

	regionLabelModal, _ := makeLabelModal(v.availableLabels, RegionLabelModal)
	imageLabelModal, _ := makeLabelModal(v.availableLabels, ImageLabelModal)

	pb.SetContent(
		Group([]Node{
			Raw(*regionLabelModal),
			Raw(*imageLabelModal),
			Table(
				Tr(Div(Class("flex"), v.ScrollerView.Render(v.scrollerButtons), v.ShapeSelector())),
				Tr(Td(Table(
					Tr(Td(Class("align-top"), v.ImageView.Build(*v.image)),
						Td(Class("align-top pl-2"),
							Div(Class("pb-2"), v.ImageInfosView.Build(*v.imageInfo)),
							Div(ID("annotation-list"), v.AnnotationsListView.Build(v.boxes, v.polygons, v.imageLabels, v.availableLabels)))),
				),
				))),
		}), nil)
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
