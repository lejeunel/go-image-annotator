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
	Update(*Image, string, time.Time) error
	DeleteImagesInCollection(*clc.Collection) error
	GetBase(ImageId) (*Image, error)
	GetAdjacent(*Image, FilterArgs, OrderingArgs, bool) (*Image, error)
	Count(FilterArgs) (int64, error)
	List(FilterArgs, OrderingArgs, g.PaginationParams) ([]Image, *g.PaginationMeta, error)
	ListWithChecksum(string) ([]Image, error)

	AssignToCollection(*Image, *clc.Collection) error
	RemoveImageFromCollection(*Image) error
	ImageIsInCollection(*Image, *clc.Collection) (bool, error)

	AssignCamera(*loc.Camera, *Image) error
	UnassignCamera(ImageId) error
}
