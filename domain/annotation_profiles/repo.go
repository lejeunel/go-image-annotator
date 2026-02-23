package annotation_profile

import (
	lbl "datahub/domain/labels"
	g "datahub/generic"
)

type AnnotationProfileRepo interface {
	Save(*AnnotationProfile) error
	Delete(*AnnotationProfile) error
	Find(AnnotationProfileId) (*AnnotationProfile, error)
	List(g.PaginationParams) ([]AnnotationProfile, *g.PaginationMeta, error)
	FindByName(string) (*AnnotationProfile, error)
	AddLabel(*AnnotationProfile, *lbl.Label) error
	RemoveLabel(*AnnotationProfile, *lbl.Label) error
	ClearLabels(*AnnotationProfile) error
	GetLabelIds(*AnnotationProfile) ([]lbl.LabelId, error)
	Rename(AnnotationProfileId, string) error
}
