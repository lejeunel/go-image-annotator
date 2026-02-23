package locationpicker

import (
	"context"
	a "datahub/app/authorizer"
	im "datahub/domain/images"
	loc "datahub/domain/locations"
	"fmt"
	"log/slog"
)

func NewLocationPicker(im *im.Service, l *loc.Service,
	auth *a.Authorizer, logger *slog.Logger) *LocationPicker {
	return &LocationPicker{
		Images:     im,
		Locations:  l,
		Authorizer: auth,
		Logger:     logger,
	}
}

type LocationPicker struct {
	Images     *im.Service
	Locations  *loc.Service
	Authorizer *a.Authorizer
	Logger     *slog.Logger
}

type LocationPickerState struct {
	ImageId           im.ImageId
	Group             string
	CanChangeLocation bool
	AvailableSites    []loc.Site
	AvailableCameras  []*loc.Camera
	Site              *loc.Site
	Camera            *loc.Camera
}

func (s *LocationPickerState) GetCameraName() string {
	if s.Camera != nil {
		return s.Camera.Name
	}
	return ""
}

func (s *LocationPickerState) GetCameraId() string {
	if s.Camera != nil {
		return s.Camera.Id.String()
	}
	return ""
}

func (s *LocationPickerState) GetSiteName() string {
	if s.Site != nil {
		return s.Site.Name
	}
	return ""
}

func (s *LocationPickerState) GetSiteId() string {
	if s.Site != nil {
		return s.Site.Id.String()
	}
	return ""
}

func (p *LocationPicker) populateAvailableCameras(ctx context.Context, site *loc.Site, state *LocationPickerState) error {
	if site == nil {
		return nil
	}
	availableCameras, err := p.Locations.ListCamerasOfSite(ctx, site)
	if err != nil {
		return fmt.Errorf("making location picker state: %w", err)
	}

	state.AvailableCameras = availableCameras
	return nil
}

func (p *LocationPicker) populateAvailableSites(ctx context.Context, state *LocationPickerState) error {
	availableSites, err := p.Locations.GetAllSites(ctx)
	if err != nil {
		return fmt.Errorf("making location picker state: %w", err)
	}

	state.AvailableSites = availableSites
	return nil
}

func (p *LocationPicker) setCurrentSite(image *im.BaseImage, state *LocationPickerState) error {
	if image.Camera != nil {
		state.Site = image.Camera.Site
	}
	return nil
}

func (p *LocationPicker) setCurrentCamera(image *im.BaseImage, state *LocationPickerState) error {
	if image.Camera != nil {
		state.Camera = image.Camera
	}
	return nil
}

func (p *LocationPicker) Init(ctx context.Context, image *im.BaseImage) (*LocationPickerState, error) {

	canChangeLocation := false
	if err := p.Authorizer.WantToContributeLocation(ctx, image.Group); err == nil {
		canChangeLocation = true
	}

	state := &LocationPickerState{
		ImageId:           image.Id,
		Group:             image.Group,
		CanChangeLocation: canChangeLocation,
	}
	if err := p.setCurrentSite(image, state); err != nil {
		return nil, fmt.Errorf("making location picker state: setting current site: %w", err)
	}
	if err := p.setCurrentCamera(image, state); err != nil {
		return nil, fmt.Errorf("making location picker state: setting current camera: %w", err)
	}

	if err := p.populateAvailableCameras(ctx, state.Site, state); err != nil {
		return nil, fmt.Errorf("making location picker state: populating cameras: %w", err)
	}
	if err := p.populateAvailableSites(ctx, state); err != nil {
		return nil, fmt.Errorf("making location picker state: populating sites: %w", err)
	}
	return state, nil

}

func (p *LocationPicker) SelectSite(ctx context.Context, site *loc.Site, image *im.BaseImage) (*LocationPickerState, error) {
	state, err := p.Init(ctx, image)

	if image.Camera != nil {
		if site.Id == image.Camera.Site.Id {
			return state, nil
		}

	}
	state.Site = site
	if err != nil {
		return nil, fmt.Errorf("selecting site %v: %w", site.Name, err)
	}

	if err := p.populateAvailableCameras(ctx, site, state); err != nil {
		return nil, fmt.Errorf("selecting site %v: %w", site.Name, err)
	}
	state.Camera = state.AvailableCameras[0]
	return state, nil
}

func (p *LocationPicker) SelectCamera(ctx context.Context, camera *loc.Camera, image *im.BaseImage) (*LocationPickerState, error) {
	state, err := p.Init(ctx, image)
	state.Camera = camera
	state.Site = camera.Site
	if err != nil {
		return nil, fmt.Errorf("selecting camera %v: %w", camera.Name, err)
	}
	if err := p.populateAvailableCameras(ctx, state.Site, state); err != nil {
		return nil, fmt.Errorf("selecting site %v: %w", state.GetSiteName(), err)
	}
	return state, nil
}
