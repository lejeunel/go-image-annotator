package locations

import (
	g "datahub/generic"
)

type LocationRepo interface {
	CreateSite(*Site) error
	CreateCamera(*Camera) error

	FindCamera(CameraId) (*Camera, error)
	FindCameraByName(*Site, string) (*Camera, error)

	FindSite(SiteId) (*Site, error)
	FindSiteByName(string) (*Site, error)
	ListCamerasOfSite(*Site) ([]*Camera, error)

	DeleteSite(SiteId) error
	DeleteCamera(CameraId) error

	UpdateSite(*Site) error
	UpdateCamera(CameraId, CameraUpdatables) (*Camera, error)

	List(FilterArgs, OrderingArgs, g.PaginationParams) ([]Site, *g.PaginationMeta, error)
	NumSites(FilterArgs) (int64, error)
}
