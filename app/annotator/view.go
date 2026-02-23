package annotator

import (
	"bytes"
	locpck "datahub/app/locationpicker"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	g "datahub/generic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	gp "maragu.dev/gomponents"
	gh "maragu.dev/gomponents/html"
	"net/http"
	"text/template"
)

type Bounds struct {
	MinX float64 `json:"minX"`
	MinY float64 `json:"minY"`
	MaxX float64 `json:"maxX"`
	MaxY float64 `json:"maxY"`
}

type Geometry struct {
	XTopLeft float64 `json:"x"`
	YTopLeft float64 `json:"y"`
	W        float64 `json:"w"`
	H        float64 `json:"h"`
	Bounds   Bounds  `json:"bounds"`
}

type Selector struct {
	Type     string   `json:"type"`
	Geometry Geometry `json:"geometry"`
}

type ATarget struct {
	AnnotationId string   `json:"annotation"`
	Selector     Selector `json:"selector"`
}

type Properties struct {
	Color string `json:"color"`
	Label string `json:"label"`
}

type Annotation struct {
	Id         uuid.UUID  `json:"id"`
	Properties Properties `json:"properties"`
	Target     ATarget    `json:"target"`
}

type AnnotationSubmission struct {
	ImageId      string  `json:"image_id"`
	CollectionId string  `json:"collection_id"`
	Annotation   ATarget `json:"annotation"`
	Label        string  `json:"label"`
}

type AnnotatorViewer struct {
	GenericViewer *g.GenericViewer
}

func NewAnnotatorViewer(genericViewer *g.GenericViewer) *AnnotatorViewer {
	return &AnnotatorViewer{GenericViewer: genericViewer}
}

func (a *AnnotatorViewer) Render(as *AnnotatorState, ls *locpck.LocationPickerState, r *http.Request, w io.Writer) error {

	description := g.MakeDescriptionTable(map[string]string{
		"id":              as.ImageId.String(),
		"collection_id":   as.CollectionId.String(),
		"collection_name": as.CollectionName,
		"group":           as.Group,
		"captured_at":     as.ImageCapturedAt.Format("2006-01-02T15:04:05.000Z"),
		"created_at":      as.ImageCreatedAt.Format("2006-01-02T15:04:05.000Z"),
		"type":            as.ImageType,
		"site":            as.SiteName,
		"camera":          as.Camera,
	},
		[]string{"id", "collection_id", "collection_name", "group",
			"captured_at", "created_at", "type", "site", "camera"})
	imageAnnotationListPanel := a.makeImageAndAnnotationListPanel(as)

	labelModal, err := a.makeLabelModal(as)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return err
	}

	javascript, err := a.makeJavascript(as)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return err
	}

	locationPicker := a.makeLocationPanel(as.Request, ls)

	body := gp.Group([]gp.Node{
		gp.Raw(labelModal),
		gh.Script(gp.Raw(javascript)),
		description,
		imageAnnotationListPanel,
		locationPicker,
	})
	head := a.makeHead()

	page := a.GenericViewer.BasePage(r.Context(), "Image", head, body)
	if err := page.Render(w); err != nil {
		fmt.Fprint(w, err.Error())
		return err
	}
	return nil

}

func (a *AnnotatorViewer) makeLabelSelectorItems(s *AnnotatorState) ([]g.SelectorItem, error) {
	if len(s.AvailableLabels) == 0 {
		return []g.SelectorItem{{Name: "NONE", Value: "NONE"}}, nil
	}

	var entries []g.SelectorItem
	for _, l := range s.AvailableLabels {
		entries = append(entries, g.SelectorItem{Name: l, Value: l})
	}

	return entries, nil

}

func (a *AnnotatorViewer) makeLocationPanel(r AnnotatorRequest, s *locpck.LocationPickerState) gp.Node {

	var siteItems []g.SelectorItem
	for _, s := range s.AvailableSites {
		siteItems = append(siteItems, g.SelectorItem{Name: s.Name, Value: s.Id.String()})
	}

	var cameraItems []g.SelectorItem
	for _, c := range s.AvailableCameras {
		cameraItems = append(cameraItems, g.SelectorItem{Name: c.Name, Value: c.Id.String()})
	}

	currentCamera := g.SelectorItem{Name: s.GetCameraName(), Value: s.GetCameraId()}
	currentSite := g.SelectorItem{Name: s.GetSiteName(), Value: s.GetSiteId()}

	submitButton := gh.A(gh.Href(fmt.Sprintf("%v&camera_id=%v",
		r.String("/ui/submit-location"), currentCamera.Value)),
		gh.Class(PrimaryButtonStyle),
		gp.Text("Submit"))

	clearCameraButton := gh.A(gh.Href(r.String("/ui/clear-camera")),
		gh.Class(DangerButtonStyle),
		gp.Text("Clear"))

	sitePOSTUrl := fmt.Sprintf("htmx.ajax('POST', '%v&site_id='+selectedItem.value, '#location-panel')",
		r.String("/ui/set-site"))
	siteSelector := g.AlpineSelector(siteItems, currentSite, sitePOSTUrl)

	cameraPOSTUrl := fmt.Sprintf("htmx.ajax('POST', '%v&camera_id='+selectedItem.value, '#location-panel')",
		r.String("/ui/set-camera"))
	cameraSelector := g.AlpineSelector(cameraItems, currentCamera, cameraPOSTUrl)

	fullPanel := gh.Div(gh.ID("location-panel"),
		gh.Table(
			gh.Tr(gh.Td(gh.Div(gh.Class("font-bold"), gp.Text("Location:")))),
			gh.Tr(gh.Td(siteSelector), gh.Td(cameraSelector), gh.Td(submitButton), gh.Td(clearCameraButton)),
		),
	)
	return fullPanel

}

func (a *AnnotatorViewer) makeAnnotationListPanel(s *AnnotatorState) gp.Node {

	labels, err := a.makeLabelSelectorItems(s)
	if err != nil {
		return gp.Text("error fetching labels")
	}

	var rows []gp.Node
	for _, box := range s.BoundingBoxes {
		labelSelector, _ := a.makeLabelSelector(s.Request,
			box.Id,
			labels,
			g.SelectorItem{Name: box.Label,
				Value: box.Label})
		rows = append(rows,
			gh.Tr(
				gh.Td(gp.If(s.CanAnnotate,
					gp.Raw(
						fmt.Sprintf(`
									<div x-data="">
										<div x-on:click="editAnnotation('%v')">
											<a href="#image-annotation-panel" class='text-black'>%v</a>
										</div>
									</div>
								`, box.Id,
							penIcon)),
				),
				),
				gh.Td(gp.Raw(makeColoredCircle(box.Color))),
				gh.Td(gp.Raw(labelSelector)),
				gh.Td(gp.If(s.CanAnnotate,
					gp.Raw(
						fmt.Sprintf(`
									<div x-data="">
										<div x-on:click="deleteAnnotation('%v')">
											<a href="#" class='text-black'>%v</a>
										</div>
									</div>
								`, box.Id,
							TrashIcon)),
				),
				),
				gh.Td(gh.Table(gh.Td(gh.Tr(gh.Class("text-sm"), gh.Td(gp.Text(box.Author))),
					gh.Tr(gh.Class("text-sm"), gh.Td(gp.Text(box.Date)))))),
				gh.Td(gp.Text(box.Id)),
			),
		)
	}

	if len(rows) > 0 {
		return gh.Table(gh.Class("border-separate border-spacing-2"),
			gh.Tr(gh.Th(), gh.Th(), gh.Th(gp.Text("label")), gh.Th(),
				gh.Th(gp.Text("details")), gh.Th(gp.Text("ID"))),
			gp.Group(rows))

	}

	return gh.Div()

}

func (a *AnnotatorViewer) makeImagePanel(s *AnnotatorState) gp.Node {
	b64str := base64.StdEncoding.EncodeToString(s.ImageData)
	annotation := a.makeAnnotionsJSONPayload(s)

	return gh.Div(gh.ID("image-panel"),
		gp.Group([]gp.Node{gh.Img(gh.ID("image"), gh.Src(fmt.Sprintf("data:%v;base64,%v",
			s.ImageMIMEType,
			b64str))),
			gh.Script(gp.Raw(annotation)),
		}))
}

func (a *AnnotatorViewer) makeImageAndAnnotationListPanel(s *AnnotatorState) gp.Node {
	annotationListPanel := a.makeAnnotationListPanel(s)
	prevNextPanel := a.makeNextPrevPanel(s.PrevImage, s.NextImage, s.Request)
	imagePanel := a.makeImagePanel(s)
	return gh.Div(gh.ID("image-annotation-panel"),
		gh.Br(),
		prevNextPanel,
		gh.Table(gh.Tr(gh.Td(imagePanel, gh.Class("align-top")),
			gh.Td(gh.Class("align-top"), gh.Div(gh.ID("annotation-panel"), gh.Class("px-2"), annotationListPanel)))),
	)
}

func (a *AnnotatorViewer) makeAnnotionsJSONPayload(s *AnnotatorState) string {

	var annotations []Annotation
	for i, box := range s.BoundingBoxes {
		x, y, w, h := box.Xc, box.Yc, box.Width, box.Height
		x -= w / 2
		y -= h / 2
		id, _ := uuid.Parse(box.Id)

		annotations = append(annotations, Annotation{
			Properties: Properties{Color: Palette[i%(len(Palette)-1)],
				Label: box.Label},
			Id: id,
			Target: ATarget{AnnotationId: box.Id,
				Selector: Selector{Type: "RECTANGLE",
					Geometry: Geometry{XTopLeft: x, YTopLeft: y, W: w, H: h,
						Bounds: Bounds{MinX: x, MinY: y, MaxX: x + w, MaxY: y + h}}}}})

	}

	jsonAnnotation, err := json.Marshal(annotations)
	if err != nil {
		return fmt.Sprintf("Failed marshalling annotation: %v",
			err.Error())
	}
	return string(jsonAnnotation[:])
}

func (a *AnnotatorViewer) makeJavascript(s *AnnotatorState) (string, error) {
	tAnnot, err := template.New("annotator").ParseFS(templatesFiles, "templates/annotator.js")
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	data := struct {
		ImageId          string
		CollectionId     string
		Annotations      string
		EnableAnnotation bool
		OriginType       string
		OriginId         string
		Ordering         string
		Descending       bool
	}{s.ImageId.String(), s.CollectionId.String(), a.makeAnnotionsJSONPayload(s),
		s.CanAnnotate, s.Request.OriginEntity, s.Request.OriginId,
		s.Request.OrderingField, s.Request.OrderingDesc}

	err = tAnnot.ExecuteTemplate(buf, "annotator", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}
func (a *AnnotatorViewer) makeOrderingSelector(r AnnotatorRequest) (string, error) {

	currentItem := g.SelectorItem{Name: r.BuildKey(), Value: r.String("image")}
	var selectableItems []g.SelectorItem
	for k, v := range r.BuildOrderingURLs("image") {
		selectableItems = append(selectableItems, g.SelectorItem{Name: k, Value: v, Disabled: false})
	}
	command := "htmx.ajax('GET', selectedItem.value)"
	selector := g.AlpineSelector(selectableItems, currentItem, command)
	buf := bytes.NewBufferString("")
	selector.Render(buf)
	return buf.String(), nil
}

func (a *AnnotatorViewer) makeLabelSelector(r AnnotatorRequest, annotationId string,
	items []g.SelectorItem, currentItem g.SelectorItem) (string, error) {

	url := r.String("ui/set-label")
	POSTUrl := fmt.Sprintf("htmx.ajax('POST', '%v&annotation_id=%v&label='+selectedItem.value, '#annotation-panel')",
		url, annotationId)
	selector := g.AlpineSelector(items, currentItem, POSTUrl)
	buf := bytes.NewBufferString("")
	selector.Render(buf)
	return buf.String(), nil
}
func (a *AnnotatorViewer) makeLabelModal(s *AnnotatorState) (string, error) {
	tLabelModal, err := template.New("labelModal").ParseFS(templatesFiles,
		"templates/label_selector.html")
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	data := struct {
		Labels       []string
		CollectionId clc.CollectionId
		ImageId      im.ImageId
	}{s.AvailableLabels, s.CollectionId, s.ImageId}
	err = tLabelModal.ExecuteTemplate(buf, "labelModal",
		data)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (a *AnnotatorViewer) makeNextPrevPanel(prevImage *im.Image, nextImage *im.Image,
	r AnnotatorRequest) gp.Node {
	var nextButton, prevButton gp.Node

	if nextImage != nil {
		nextRequest := r
		nextRequest.CollectionId = nextImage.CollectionId
		nextRequest.ImageId = nextImage.Id
		nextButton = gh.A(gh.Href(nextRequest.String("image")+"#image-annotation-panel"),
			gh.Class(PrimaryButtonStyle),
			gp.Text("Next"))
	} else {
		nextButton = gh.A(gh.Href("#"),
			gh.Class(SecondaryButtonStyle),
			gp.Text("Next"))
	}

	if prevImage != nil {
		prevRequest := r
		prevRequest.CollectionId = prevImage.CollectionId
		prevRequest.ImageId = prevImage.Id
		prevButton = gh.A(gh.Href(prevRequest.String("image")+"#image-annotation-panel"),
			gh.Class(PrimaryButtonStyle),
			gp.Text("Previous"))

	} else {
		prevButton = gh.A(gh.Href("#"),
			gh.Class(SecondaryButtonStyle),
			gp.Text("Previous"))
	}
	orderingModal, err := a.makeOrderingSelector(r)
	if err != nil {
		return gp.Text(err.Error())
	}

	scrollerInfo := gh.Div(gh.Class("py-4"),
		gh.Table(
			gh.Tr(gh.Td(gh.Div(gh.Class("font-bold"), gp.Text("criteria:"))),
				gh.Td(gp.Text(r.Filter.String()))),
			gh.Tr(gh.Td(gh.Div(gh.Class("font-bold"), gp.Text("ordering:"))),
				gh.Td(gp.Raw(orderingModal)))),
	)

	return gh.Div(gh.Class("py-4"), prevButton, nextButton, scrollerInfo)

}
func (a *AnnotatorViewer) makeHead() gp.Node {
	return gp.Group([]gp.Node{
		gh.Link(gh.Href("/static/styles.css"), gh.Rel("stylesheet")),
		gh.Script(gh.Defer(), gh.Src("/static/alpine.focus.js")),
		gh.Script(gh.Defer(), gh.Src("/static/alpine.js")),
		gh.Script(gh.Defer(), gh.Src("/static/annotorious.js")),
		gh.Link(gh.Href("/static/annotorious.css"), gh.Rel("stylesheet")),
		gh.Script(gh.Defer(), gh.Src("/static/htmx.js")),
	})
}

func (a *AnnotatorViewer) makeDescription(s *AnnotatorState) gp.Node {
	return g.MakeDescriptionTable(map[string]string{
		"id":              s.ImageId.String(),
		"collection_id":   s.CollectionId.String(),
		"collection_name": s.CollectionName,
		"site":            s.SiteName,
		"camera":          s.Camera,
		"created_at":      s.ImageCreatedAt.Format("2006-01-02 / 15:04"),
		"captured_at":     s.ImageCapturedAt.Format("2006-01-02 / 15:04"),
		"type":            s.ImageType},
		[]string{"id", "collection_name", "collection_id",
			"captured_at", "created_at", "type",
			"site", "camera"})

}
