package locations

import (
	"context"
	e "datahub/errors"
	g "datahub/generic"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

type SiteInputCreate struct {
	Body struct {
		Name  string `doc:"Name" json:"name"`
		Group string `json:"group" doc:"Group" required:"true"`
	}
}

type PathId struct {
	Id string `path:"id"`
}

type GetSiteByNameRequest struct {
	Name string `path:"name" doc:"name of site"`
}

type SiteUpdate struct {
	Id   string `path:"id"`
	Body struct {
		Name  string `doc:"Name" json:"name"`
		Group string `doc:"Group" json:"group_name"`
	}
}

type CameraUpdatables struct {
	Name        string `json:"name" doc:"Name"`
	SiteName    string `json:"site" doc:"Site Name"`
	Transmitter string `json:"transmitter" required:"false" doc:"Transmitter Name"`
}

type CameraDelete struct {
	Id string `path:"id"`
}

type CameraRequestBody struct {
	Name string `doc:"Name" json:"name"`
}

type CameraUpdateRequest struct {
	Id   string `path:"id"`
	Body CameraUpdatables
}

type CameraPatchRequest struct {
	Id   string `path:"id"`
	Body g.JSONPatches
}

type CameraCreateRequest struct {
	Body struct {
		Name        string `json:"name" doc:"Camera name"`
		SiteName    string `json:"site" doc:"Site name"`
		Transmitter string `json:"transmitter" doc:"Transmitter name" required:"false"`
	}
}

type SiteResponseBody struct {
	Id      string                   `json:"id"`
	Name    string                   `json:"name"`
	Group   string                   `json:"group"`
	Cameras []CameraOfSiteOutputBody `json:"cameras"`
}

type SiteResponse struct {
	Body SiteResponseBody
}

type CameraOutputBody struct {
	Id          string `doc:"Id" json:"id"`
	Name        string `doc:"Name" json:"name"`
	SiteName    string `doc:"Site Name" json:"site"`
	Transmitter string `doc:"Transmitter Name" json:"transmitter"`
}
type CameraOutput struct {
	Body CameraOutputBody
}

func NewCameraOutput(camera *Camera) *CameraOutput {

	return &CameraOutput{Body: CameraOutputBody{Id: camera.Id.String(),
		Name:        camera.Name,
		SiteName:    camera.Site.Name,
		Transmitter: camera.Transmitter}}
}

type CameraOfSiteOutputBody struct {
	Id   string `doc:"Id" json:"id"`
	Name string `doc:"Name" json:"name"`
}

func makeSiteResponse(s Site, cameras []*Camera) *SiteResponse {
	var camerasOutput []CameraOfSiteOutputBody
	for _, c := range cameras {
		camerasOutput = append(camerasOutput, CameraOfSiteOutputBody{Name: c.Name, Id: c.Id.String()})
	}
	return &SiteResponse{Body: SiteResponseBody{Id: s.Id.String(),
		Name: s.Name, Group: s.Group, Cameras: camerasOutput}}
}

type LocationHTTPController struct {
	LocationService *Service
}

func (h *LocationHTTPController) GetSite(ctx context.Context, input *GetSiteByNameRequest) (*SiteResponse, error) {
	site, err := h.LocationService.FindSiteByName(ctx, input.Name)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error(), err)
	}
	cameras, err := h.LocationService.ListCamerasOfSite(ctx, site)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error(), err)
	}

	return makeSiteResponse(*site, cameras), nil
}

func (h *LocationHTTPController) DeleteSite(ctx context.Context, input *PathId) (*struct{}, error) {
	siteId, err := NewSiteIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	if err := h.LocationService.DeleteSite(ctx, *siteId); err != nil {
		return nil, huma.Error404NotFound(err.Error(), err)

	}

	return nil, nil
}

func (h *LocationHTTPController) CreateSite(ctx context.Context, input *SiteInputCreate) (*SiteResponse, error) {

	site, err := NewSite(input.Body.Name, WithGroupOption(input.Body.Group))
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error(), err)
	}
	if err := h.LocationService.SaveSite(ctx, site); err != nil {

		return nil, huma.Error400BadRequest(err.Error(), err)
	}
	cameras, err := h.LocationService.ListCamerasOfSite(ctx, site)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error(), err)
	}

	return makeSiteResponse(*site, cameras), nil
}

func (h *LocationHTTPController) PutSite(ctx context.Context, input *SiteUpdate) (*SiteResponse, error) {
	siteId, err := NewSiteIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	site, err := h.LocationService.FindSite(ctx, *siteId)
	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}

	site.Name = input.Body.Name
	site.Group = input.Body.Group

	site, err = h.LocationService.UpdateSite(ctx, site)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}

	return &SiteResponse{Body: SiteResponseBody{Id: site.Id.String(), Name: site.Name}}, nil
}

func (h *LocationHTTPController) AddCamera(ctx context.Context, input *CameraCreateRequest) (*CameraOutput, error) {
	site, err := h.LocationService.FindSiteByName(ctx, input.Body.SiteName)
	if err != nil {
		return nil, fmt.Errorf("adding camera: fetching site by name (%v): %w", input.Body.SiteName, err)
	}

	camera, err := NewCamera(input.Body.Name, site, WithTransmitter(input.Body.Transmitter))
	if err = h.LocationService.SaveCamera(ctx, camera); err != nil {
		return nil, fmt.Errorf("adding camera: %w", err)
	}

	return NewCameraOutput(camera), nil
}

func (h *LocationHTTPController) DeleteCamera(ctx context.Context, input *CameraDelete) (*struct{}, error) {
	camId, err := NewCameraIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	if err := h.LocationService.DeleteCamera(ctx, *camId); err != nil {
		return nil, fmt.Errorf("deleting camera: %w:", err)
	}
	return nil, nil
}

func (h *LocationHTTPController) PutCamera(ctx context.Context, input *CameraUpdateRequest) (*CameraOutput, error) {
	camId, err := NewCameraIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	camera, err := h.LocationService.UpdateCamera(ctx, *camId, input.Body)
	if err != nil {
		return nil, fmt.Errorf("updating camera: %w", err)
	}

	return NewCameraOutput(camera), nil
}

func (h *LocationHTTPController) PatchCamera(ctx context.Context, input *CameraPatchRequest) (*CameraOutput, error) {
	camId, err := NewCameraIdFromString(input.Id)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	camera, err := h.LocationService.PatchCamera(ctx, *camId, input.Body)
	if err != nil {
		return nil, fmt.Errorf("patching camera: %w", err)
	}

	return NewCameraOutput(camera), nil
}
