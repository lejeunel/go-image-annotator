package services

import (
	"context"
	"github.com/google/uuid"
	e "go-image-annotator/errors"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
)

type LabelService struct {
	LabelRepo       r.LabelRepo
	ImageRepo       r.ImageRepo
	MaxPageSize     int
	DefaultPageSize int
}

func (s *LabelService) Create(ctx context.Context, label *m.Label) (*m.Label, error) {

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

func (s *LabelService) GetOne(ctx context.Context, id string) (*m.Label, error) {

	label, err := s.LabelRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return label, nil

}

func (s *LabelService) Delete(ctx context.Context, label *m.Label) error {

	numImages, err := s.LabelRepo.NumImagesWithLabel(ctx, label)
	if err != nil {
		return err
	}

	if numImages > 0 {
		return e.ErrForbiddenDeletingDependency{ParentEntity: "image", ParentId: label.Id.String(), ChildEntity: "label"}
	}

	return s.LabelRepo.Delete(ctx, label)

}
