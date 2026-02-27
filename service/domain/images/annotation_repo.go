package images

import (
	clc "datahub/domain/collections"
	lbl "datahub/domain/labels"
)

type AnnotationRepo interface {
	GetAnnotationIdsOfImage(*Image) ([]string, error)
	GetAnnotationById(string) (*Annotation, error)
	ApplyLabelToImage(*lbl.Label, *Image, string) error
	DeleteAnnotation(*Annotation) error
	DeleteAllAnnotations(*clc.Collection) error
	UpdateBoundingBox(*BoundingBox, *Image) error
	ApplyBoundingBox(*BoundingBox, *Image) error
	GetBoundingBoxesOfImage(*Image) ([]*BoundingBox, error)
	UpdateAnnotationLabel(annotationId string, labelId string) error
}
