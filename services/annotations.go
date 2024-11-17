package services

import (
	"context"
	"github.com/google/uuid"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
)

type AnnotationService struct {
	LabelRepo       r.LabelRepo
	ImageRepo       r.ImageRepo
	MaxPageSize     int
	DefaultPageSize int
}

func (s *AnnotationService) Create(ctx context.Context, label *m.Label) (*m.Label, error) {

	if err := label.Validate(); err != nil {
		return nil, err
	}

	label.Id = uuid.New()

	label, err := s.LabelRepo.Create(ctx, label)
	if err != nil {
		return nil, err
	}

	return label, nil

}

func (s *AnnotationService) GetOne(ctx context.Context, id string) (*m.Label, error) {

	label, err := s.LabelRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return label, nil

}

func (s *AnnotationService) Delete(ctx context.Context, label *m.Label) error {

	numImages, err := s.LabelRepo.NumImagesWithLabel(ctx, label)
	if err != nil {
		return err
	}

	if numImages > 0 {
		return e.ErrForbiddenDeletingDependency{ParentEntity: "image", ParentId: label.Id.String(), ChildEntity: "label"}
	}

	return s.LabelRepo.Delete(ctx, label)

}

func (s *AnnotationService) ApplyLabelToImage(ctx context.Context, label *m.Label, image *m.Image) (*m.Image, error) {
	if err := s.LabelRepo.ApplyLabelToImage(ctx, label, image); err != nil {
		return nil, err
	}

	labels, err := s.LabelRepo.GetLabelsOfImage(ctx, image)
	if err != nil {
		return nil, err
	}

	image.Labels = labels

	return image, nil

}
func (s *AnnotationService) ApplyPolygonToImage(ctx context.Context, polygon *m.Polygon, image *m.Image) (*m.Image, error) {
	image.Polygons = append(image.Polygons, polygon)

	if err := s.LabelRepo.ApplyPolygonToImage(ctx, polygon, image); err != nil {
		return nil, err
	}
	return image, nil
}
