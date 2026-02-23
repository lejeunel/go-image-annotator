package collections

import (
	pro "datahub/domain/annotation_profiles"
	g "datahub/generic"
	"time"
)

type CollectionRepo interface {
	Create(*Collection) error
	Find(CollectionId) (*Collection, error)
	GetByName(string) (*Collection, error)
	Delete(*Collection) error
	Update(CollectionId, CollectionUpdatables) error
	Touch(CollectionId, time.Time) error
	List(OrderingArgs, g.PaginationParams) ([]Collection, *g.PaginationMeta, error)
	Count() (int64, error)

	AssignProfile(*Collection, *pro.AnnotationProfile) error
	UnassignProfile(*Collection) error
}
