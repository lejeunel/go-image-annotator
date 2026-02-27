package images

import (
	clc "datahub/domain/collections"
	loc "datahub/domain/locations"
	g "datahub/generic"
	"time"
)

type ImageRepo interface {
	Create(*Image) (*Image, error)
	Delete(*Image) error
	Update(ImageId, string, time.Time) error
	DeleteImagesInCollection(*clc.Collection) error
	GetBase(ImageId) (*BaseImage, error)
	GetAdjacent(*Image, FilterArgs, OrderingArgs, bool) (*BaseImage, error)
	Count(FilterArgs) (int64, error)
	List(FilterArgs, OrderingArgs, g.PaginationParams) ([]BaseImage, *g.PaginationMeta, error)
	ListWithChecksum(string) ([]BaseImage, error)

	AssignToCollection(*Image, *clc.Collection) error
	RemoveImageFromCollection(*Image) error
	ImageIsInCollection(*Image, *clc.Collection) (bool, error)

	AssignCamera(loc.CameraId, ImageId) error
	UnassignCamera(ImageId) error
}
