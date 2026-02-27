package locations

import (
	"context"
	au "datahub/app/authorizer"
	e "datahub/errors"
	g "datahub/generic"
	"encoding/json"
	"fmt"
	clk "github.com/jonboulle/clockwork"
	"log/slog"
)

type Service struct {
	LocationRepo    LocationRepo
	MaxPageSize     int
	DefaultPageSize int
	Authorizer      *au.Authorizer
	Logger          *slog.Logger
	Clock           clk.Clock
}

func NewLocationService(repo LocationRepo, maxPageSize int, defaultPageSize int, auth *au.Authorizer,
	logger *slog.Logger, clock clk.Clock) *Service {

	return &Service{
		LocationRepo:    repo,
		MaxPageSize:     maxPageSize,
		DefaultPageSize: defaultPageSize,
		Authorizer:      auth,
		Clock:           clock,
		Logger:          logger,
	}
}

func (s *Service) SaveSite(ctx context.Context, site *Site) error {
	if err := site.Validate(); err != nil {
		return err
	}
	err := s.Authorizer.WantToContributeLocation(ctx, site.Group)
	if err != nil {
		return fmt.Errorf("creating site: %w", err)
	}
	site.CreatedAt = s.Clock.Now()
	site.UpdatedAt = s.Clock.Now()

	if err := s.LocationRepo.CreateSite(site); err != nil {
		return err
	}
	return nil
}

func (s *Service) SaveCamera(ctx context.Context, camera *Camera) error {
	if err := s.Authorizer.WantToContributeLocation(ctx, camera.Group); err != nil {
		return fmt.Errorf("creating camera: %w", err)
	}

	camera.CreatedAt = s.Clock.Now()
	camera.UpdatedAt = s.Clock.Now()

	if err := s.LocationRepo.CreateCamera(camera); err != nil {
		return err
	}
	site, err := s.FindSite(ctx, camera.Site.Id)
	if err != nil {
		return fmt.Errorf("fetching site of camera: %w", err)
	}

	camera.Site.Name = site.Name
	return nil
}

func (s *Service) ListCamerasOfSite(ctx context.Context, site *Site) ([]*Camera, error) {
	cameras, err := s.LocationRepo.ListCamerasOfSite(site)
	if err != nil {
		return nil, err
	}
	return cameras, nil
}

func (s *Service) FindCamera(ctx context.Context, id CameraId) (*Camera, error) {
	camera, err := s.LocationRepo.FindCamera(id)
	if err != nil {
		return nil, fmt.Errorf("fetching camera by id %v: %w", id, err)
	}

	return camera, nil
}

func (s *Service) FindSite(ctx context.Context, id SiteId) (*Site, error) {
	site, err := s.LocationRepo.FindSite(id)
	if err != nil {
		return nil, fmt.Errorf("fetching site by id %v: %w", id, err)
	}
	return site, nil
}

func (s *Service) FindSiteByName(ctx context.Context, name string) (*Site, error) {
	site, err := s.LocationRepo.FindSiteByName(name)
	if err != nil {
		return nil, fmt.Errorf("fetching site by name %v: %w", name, err)
	}

	return site, nil
}

func (s *Service) GetAllSites(ctx context.Context) ([]Site, error) {
	numSites, err := s.LocationRepo.NumSites(FilterArgs{})
	if err != nil {
		return nil, fmt.Errorf("fetching all sites: %w", err)
	}

	sites, _, err := s.List(ctx, FilterArgs{}, SiteAlphabeticalOrdering,
		g.PaginationParams{Page: 1, PageSize: int(numSites)})
	if err != nil {
		return nil, fmt.Errorf("fetching all sites: %w", err)
	}

	return sites, nil
}

func (s *Service) List(
	ctx context.Context,
	filters FilterArgs,
	ordering OrderingArgs,
	pagination g.PaginationParams) ([]Site, *g.PaginationMeta, error) {
	sites, paginationMeta, err := s.LocationRepo.List(filters, ordering, pagination)
	if err != nil {
		return nil, nil, err
	}

	return sites, paginationMeta, nil

}

func (s *Service) DeleteSite(ctx context.Context, id SiteId) error {
	site, err := s.FindSite(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching site by id %v: %w", id, err)
	}

	if err := s.Authorizer.WantToContributeLocation(ctx, site.Group); err != nil {
		return fmt.Errorf("deleting site: %w", err)
	}
	return s.LocationRepo.DeleteSite(id)
}

func (s *Service) DeleteCamera(ctx context.Context, id CameraId) error {
	camera, err := s.FindCamera(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching camera by id %v: %w", id, err)
	}
	if err := s.Authorizer.WantToContributeLocation(ctx, camera.Group); err != nil {
		return fmt.Errorf("deleting camera: %w", err)
	}
	return s.LocationRepo.DeleteCamera(id)
}

func (s *Service) siteMustNotHaveCameraWithName(ctx context.Context, siteName string, cameraId CameraId) error {
	baseErrMsg := fmt.Sprintf("checking whether site %v has a camera with id %v", siteName, cameraId)
	currentCamera, err := s.FindCamera(ctx, cameraId)
	if err != nil {
		return fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	if siteName != currentCamera.Site.Name {
		newSite, err := s.FindSiteByName(ctx, siteName)
		if err != nil {
			return fmt.Errorf("%v: %w", baseErrMsg, err)
		}
		cameras, err := s.ListCamerasOfSite(ctx, newSite)
		for _, c := range cameras {
			if c.Name == currentCamera.Name {
				return fmt.Errorf("%v: camera with name %v already exists: %w",
					baseErrMsg, currentCamera.Name, e.ErrDuplication)
			}
		}

	}
	return nil

}
func (s *Service) FindCameraByName(ctx context.Context, site *Site, cameraName string) (*Camera, error) {
	return s.LocationRepo.FindCameraByName(site, cameraName)
}

func (s *Service) UpdateCamera(ctx context.Context, id CameraId, update CameraUpdatables) (*Camera, error) {
	camera, err := s.FindCamera(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fetching camera by id %v: %w", id, err)
	}
	if err := s.Authorizer.WantToContributeLocation(ctx, camera.Group); err != nil {
		return nil, fmt.Errorf("updating camera: %w", err)
	}

	if err := s.siteMustNotHaveCameraWithName(ctx, update.SiteName, id); err != nil {
		return nil, fmt.Errorf("updating camera: %w", err)
	}

	return s.LocationRepo.UpdateCamera(id, update)
}

func (s *Service) PatchCamera(ctx context.Context, id CameraId, patches g.JSONPatches) (*Camera, error) {
	baseErrMsg := "patching camera"
	camera, err := s.FindCamera(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, id, err)
	}
	if err := s.Authorizer.WantToContributeLocation(ctx, camera.Group); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}
	originalJSONBytes, err := json.Marshal(
		&CameraUpdatables{Name: camera.Name,
			SiteName:    camera.Site.Name,
			Transmitter: camera.Transmitter})
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	modifiedBytes, err := patches.Apply(originalJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	var modified CameraUpdatables
	if err := json.Unmarshal(modifiedBytes, &modified); err != nil {
		return nil, fmt.Errorf("%v: %w", baseErrMsg, err)
	}

	return s.UpdateCamera(ctx, id, modified)
}

func (s *Service) UpdateSite(ctx context.Context, site *Site) (*Site, error) {
	err := s.Authorizer.WantToContributeLocation(ctx, site.Group)
	if err != nil {
		return nil, fmt.Errorf("updating site: %w", err)
	}

	if err := site.Validate(); err != nil {
		return nil, fmt.Errorf("validating collection prior to update: %v", err)
	}

	col, err := s.FindSite(ctx, site.Id)
	if err != nil {
		return nil, fmt.Errorf("getting site %v prior to update: %w", col, err)
	}

	err = s.LocationRepo.UpdateSite(site)
	if err != nil {
		return nil, err
	}

	return site, nil
}
