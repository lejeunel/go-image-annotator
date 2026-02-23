package annotator

import (
	"context"
	locpck "datahub/app/locationpicker"
	clc "datahub/domain/collections"
	im "datahub/domain/images"
	lbl "datahub/domain/labels"
	loc "datahub/domain/locations"
	g "datahub/generic"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type BoundingBox struct {
	Id     string
	Xc     float64
	Yc     float64
	Width  float64
	Height float64
	Angle  float64
	Author string
	Label  string
	Date   string
	Color  string
}

type AnnotatorRequest struct {
	ImageId       im.ImageId
	CollectionId  clc.CollectionId
	Filter        im.FilterArgs
	Ordering      im.OrderingArgs
	OriginEntity  string
	OriginId      string
	OrderingField string
	OrderingDesc  bool
	SiteId        *loc.SiteId
	CameraId      *loc.CameraId
}

func NewAnnotatorRequestFromURL(url url.Values) (*AnnotatorRequest, error) {

	collectionId, err := clc.NewCollectionIdFromString(url.Get("collection_id"))
	if err != nil {
		return nil, err
	}

	imageId, err := im.NewImageIdFromString(url.Get("image_id"))
	if err != nil {
		return nil, err
	}

	annotatorRequest := &AnnotatorRequest{ImageId: *imageId, CollectionId: *collectionId}
	originEntity := url.Get("origin_entity")
	originId := url.Get("origin_id")
	ordering := url.Get("ordering")
	descending := url.Get("descending")

	annotatorRequest.OriginEntity = originEntity
	annotatorRequest.OriginId = originId

	switch originEntity {
	case "collection":

		annotatorRequest.Filter = *im.NewImageFilter(im.WithCollectionId(*collectionId))
	case "label":
		labelId, err := lbl.NewLabelIdFromString(originId)
		if err != nil {
			return nil, err
		}
		annotatorRequest.Filter = *im.NewImageFilter(im.WithLabelId(*labelId))
	case "camera":
		camId, err := loc.NewCameraIdFromString(originId)
		if err != nil {
			return nil, err
		}
		annotatorRequest.Filter = *im.NewImageFilter(im.WithCameraId(*camId))
	default:
		annotatorRequest.Filter = *im.NewImageFilter(im.WithCollectionId(*collectionId))
	}

	switch {
	case (ordering == "captured_at" && descending == "false"):
		annotatorRequest.Ordering = *im.NewAscendingImageCapturedOrder()
		annotatorRequest.OrderingField = "captured_at"
		annotatorRequest.OrderingDesc = false
	case (ordering == "captured_at" && descending == "true"):
		annotatorRequest.Ordering = *im.NewDescendingImageCapturedOrder()
		annotatorRequest.OrderingField = "captured_at"
		annotatorRequest.OrderingDesc = true
	case (ordering == "created_at" && descending == "false"):
		annotatorRequest.Ordering = *im.NewAscendingImageCreatedOrder()
		annotatorRequest.OrderingField = "created_at"
		annotatorRequest.OrderingDesc = false
	case (ordering == "created_at" && descending == "true"):
		annotatorRequest.Ordering = *im.NewDescendingImageCreatedOrder()
		annotatorRequest.OrderingField = "created_at"
		annotatorRequest.OrderingDesc = true
	default:
		annotatorRequest.Ordering = *im.NewAscendingImageCapturedOrder()
		annotatorRequest.OrderingField = "captured_at"
		annotatorRequest.OrderingDesc = false
	}

	return annotatorRequest, nil
}

func (r *AnnotatorRequest) BuildKey() string {

	return fmt.Sprintf("%v / desc: %v", r.OrderingField, r.OrderingDesc)
}

func (r *AnnotatorRequest) BuildOrderingURLs(prefix string) map[string]string {
	m := make(map[string]string)
	orderingFields := []string{"captured_at", "created_at"}
	orderingDesc := []bool{false, true}
	for _, f := range orderingFields {
		for _, d := range orderingDesc {
			rCopy := r
			rCopy.OrderingField = f
			rCopy.OrderingDesc = d
			m[rCopy.BuildKey()] = rCopy.String(prefix)
		}
	}
	return m
}

func (r *AnnotatorRequest) String(prefix string) string {
	return fmt.Sprintf("%v?image_id=%v&collection_id=%v&origin_entity=%v&origin_id=%v&ordering=%v&descending=%v",
		prefix, r.ImageId, r.CollectionId, r.OriginEntity, r.OriginId, r.OrderingField, r.OrderingDesc)
}

type AnnotatorController struct {
	Locations *AnnotatorLocationController
	Model     *Annotator
	View      *AnnotatorViewer
	Logger    *slog.Logger
}

func NewAnnotatorController(annotator *Annotator, viewer *AnnotatorViewer, logger *slog.Logger) *AnnotatorController {
	return &AnnotatorController{
		Locations: &AnnotatorLocationController{Annotator: annotator, Viewer: viewer, Logger: logger},
		Model:     annotator,
		View:      viewer,
		Logger:    logger,
	}
}

func (a *AnnotatorController) RegisterEndPoints(mux *http.ServeMux) {

	mux.HandleFunc("GET /ui/annotation-panel", a.AnnotationPanelHandler)
	mux.HandleFunc("POST /ui/submit-annotation", a.SubmitAnnotationHandler)
	mux.HandleFunc("POST /ui/set-label", a.SetLabelHandler)
	mux.HandleFunc("DELETE /ui/delete-annotation", a.DeleteAnnotationHandler)
	mux.HandleFunc("GET /image", a.mainHandler)

	mux.HandleFunc("POST /ui/set-site", a.Locations.SelectSiteHandler)
	mux.HandleFunc("POST /ui/set-camera", a.Locations.SelectCameraHandler)
	mux.HandleFunc("GET /ui/clear-camera", a.ClearCameraHandler)
	mux.HandleFunc("GET /ui/submit-location", a.SubmitLocationHandler)
}
func (a *AnnotatorController) convertBoundingBoxPayload(ctx context.Context, an *AnnotationSubmission) (*BoundingBox, error) {

	w, h := an.Annotation.Selector.Geometry.W, an.Annotation.Selector.Geometry.H
	xc := an.Annotation.Selector.Geometry.XTopLeft + w/2
	yc := an.Annotation.Selector.Geometry.YTopLeft + h/2

	author, err := a.Model.Authorizer.IdentityProvider.Email(ctx)
	if err != nil {
		return nil, fmt.Errorf("building  bounding box payload: %w", err)
	}

	bboxPayload := BoundingBox{
		Id:     an.Annotation.AnnotationId,
		Xc:     xc,
		Yc:     yc,
		Width:  w,
		Height: h,
		Angle:  0,
		Author: author,
		Label:  an.Label,
		Date:   time.Now().Format("2006-01-02T15:04:05.000Z")}

	return &bboxPayload, nil
}
func (a *AnnotatorController) parseAnnotation(w *http.ResponseWriter, r *http.Request) (*AnnotationSubmission, error) {
	bodyBytes, _ := io.ReadAll(r.Body)

	var an AnnotationSubmission
	err := json.Unmarshal(bodyBytes, &an)
	if err != nil {
		g.LogAndWriteError(a.Logger,
			fmt.Errorf("submitting annotation: reading json payload: %W", err), w)
		return nil, err
	}
	return &an, nil
}

func (a *AnnotatorController) mainHandler(w http.ResponseWriter, r *http.Request) {
	locPicker := locpck.NewLocationPicker(a.Model.Images,
		a.Model.Locations, a.Model.Authorizer,
		a.Logger)

	stateRequest, err := NewAnnotatorRequestFromURL(r.URL.Query())
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	annotatorState, err := a.Model.MakeState(r.Context(), *stateRequest)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	image, err := a.Model.Images.Find(r.Context(), stateRequest.ImageId, stateRequest.CollectionId,
		im.FetchMetaOnly)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	locationPickerState, err := locPicker.Init(r.Context(), image)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	if err := a.View.Render(annotatorState, locationPickerState, r, w); err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
	}

}
func (a *AnnotatorController) SubmitLocationHandler(w http.ResponseWriter, r *http.Request) {
	a.Locations.submitLocationHandler(w, r)
	a.mainHandler(w, r)
}
func (a *AnnotatorController) ClearCameraHandler(w http.ResponseWriter, r *http.Request) {
	a.Locations.ClearCameraHandler(w, r)
	a.mainHandler(w, r)
}
func (a *AnnotatorController) SetLabelHandler(w http.ResponseWriter, r *http.Request) {

	baseErrMsg := "setting label"
	stateRequest, err := NewAnnotatorRequestFromURL(r.URL.Query())
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	state, err := a.Model.MakeState(r.Context(), *stateRequest)
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

	labelName := r.URL.Query().Get("label")
	annotationId := r.URL.Query().Get("annotation_id")
	if err := a.Model.Images.Annotations.UpdateLabel(r.Context(), annotationId, labelName); err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

	state, err = a.Model.MakeState(r.Context(), *stateRequest)
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
	}
	body := a.View.makeAnnotationListPanel(state)
	if err := body.Render(w); err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}
}
func (a *AnnotatorController) AnnotationPanelHandler(w http.ResponseWriter, r *http.Request) {

	baseErrMsg := "building annotation panel"
	stateRequest, err := NewAnnotatorRequestFromURL(r.URL.Query())
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}
	state, err := a.Model.MakeState(r.Context(), *stateRequest)

	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}
	body := a.View.makeAnnotationListPanel(state)
	if err := body.Render(w); err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

}
func (a *AnnotatorController) DeleteAnnotationHandler(w http.ResponseWriter, r *http.Request) {
	baseErrMsg := "building annotation panel"
	stateRequest, err := NewAnnotatorRequestFromURL(r.URL.Query())
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}
	annotationId := r.URL.Query().Get("annotation_id")
	if err := a.Model.DeleteAnnotation(r.Context(), annotationId); err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

	state, err := a.Model.MakeState(r.Context(), *stateRequest)
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(a.View.makeAnnotionsJSONPayload(state)))

}
func (a *AnnotatorController) SubmitAnnotationHandler(w http.ResponseWriter, r *http.Request) {
	baseErrMsg := "submiting annotation"
	an, err := a.parseAnnotation(&w, r)
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

	bbox, err := a.convertBoundingBoxPayload(r.Context(), an)
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

	stateRequest, err := NewAnnotatorRequestFromURL(r.URL.Query())
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}
	state, err := a.Model.MakeState(r.Context(), *stateRequest)
	if err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}
	if err := a.Model.UpsertBoundingBox(r.Context(), state.ImageId, state.CollectionId, bbox); err != nil {
		g.LogAndWriteError(a.Logger, fmt.Errorf("%v: %w", baseErrMsg, err), &w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	state, _ = a.Model.MakeState(r.Context(), *stateRequest)

	jsonPayload := a.View.makeAnnotionsJSONPayload(state)
	w.Write([]byte(jsonPayload))
}
