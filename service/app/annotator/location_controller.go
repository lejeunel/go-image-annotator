package annotator

import (
	locpck "datahub/app/locationpicker"
	im "datahub/domain/images"
	loc "datahub/domain/locations"
	g "datahub/generic"
	"log/slog"
	"net/http"
)

type AnnotatorLocationController struct {
	Annotator *Annotator
	Viewer    *AnnotatorViewer
	Logger    *slog.Logger
}

func (a *AnnotatorLocationController) initLocationHandlers(r *http.Request) (*locpck.LocationPicker, *im.BaseImage, error) {
	locPicker := locpck.NewLocationPicker(a.Annotator.Images,
		a.Annotator.Locations, a.Annotator.Authorizer,
		a.Logger)

	imageId, err := im.NewImageIdFromString(r.URL.Query().Get("image_id"))
	if err != nil {
		return nil, nil, err
	}

	image, err := a.Annotator.Images.GetBase(r.Context(), *imageId, im.FetchMetaOnly)
	if err != nil {
		return nil, nil, err
	}
	return locPicker, image, nil
}
func (a *AnnotatorLocationController) SelectSiteHandler(w http.ResponseWriter, r *http.Request) {

	locPicker, image, err := a.initLocationHandlers(r)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	siteId, err := loc.NewSiteIdFromString(r.URL.Query().Get("site_id"))
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	site, err := a.Annotator.Locations.FindSite(r.Context(), *siteId)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	locationPickerState, err := locPicker.SelectSite(r.Context(), site, image)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	annotatorState, err := NewAnnotatorRequestFromURL(r.URL.Query())
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	a.Viewer.makeLocationPanel(*annotatorState, locationPickerState).Render(w)

}
func (a *AnnotatorLocationController) ClearCameraHandler(w http.ResponseWriter, r *http.Request) {
	imageId, err := im.NewImageIdFromString(r.URL.Query().Get("image_id"))
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	if err := a.Annotator.Images.UnassignCamera(r.Context(), *imageId); err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
}
func (a *AnnotatorLocationController) SelectCameraHandler(w http.ResponseWriter, r *http.Request) {

	locPicker, image, err := a.initLocationHandlers(r)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	cameraId, err := loc.NewCameraIdFromString(r.URL.Query().Get("camera_id"))
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	camera, err := a.Annotator.Locations.FindCamera(r.Context(), *cameraId)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	locationPickerState, err := locPicker.SelectCamera(r.Context(), camera, image)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	annotatorState, err := NewAnnotatorRequestFromURL(r.URL.Query())
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	a.Viewer.makeLocationPanel(*annotatorState, locationPickerState).Render(w)

}
func (a *AnnotatorLocationController) submitLocationHandler(w http.ResponseWriter, r *http.Request) {
	cameraId, err := loc.NewCameraIdFromString(r.URL.Query().Get("camera_id"))
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	imageId, err := im.NewImageIdFromString(r.URL.Query().Get("image_id"))
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	image, err := a.Annotator.Images.GetBase(r.Context(), *imageId, im.FetchMetaOnly)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}
	camera, err := a.Annotator.Locations.FindCamera(r.Context(), *cameraId)
	if err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

	if err := a.Annotator.Images.AssignCamera(r.Context(), camera.Id, image.Id); err != nil {
		g.LogAndWriteError(a.Logger, err, &w)
		return
	}

}
