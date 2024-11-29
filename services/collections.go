package services

import (
	"context"
	"github.com/google/uuid"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
)

type CollectionService struct {
	CollectionRepo  r.CollectionRepo
	ImageRepo       r.ImageRepo
	ImageService    ImageService
	MaxPageSize     int
	DefaultPageSize int
}

func (s *CollectionService) Create(ctx context.Context, set *m.Collection) (*m.Collection, error) {
	if err := set.Validate(); err != nil {
		return nil, err
	}

	set.Id = uuid.New()
	set, err := s.CollectionRepo.Create(ctx, set)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (s *CollectionService) GetOne(ctx context.Context, id string) (*m.Collection, error) {
	set, err := s.CollectionRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (s *CollectionService) AppendImageToSet(ctx context.Context, image *m.Image, set *m.Collection) error {
	err := s.CollectionRepo.AssignImageToSet(ctx, image, set)
	if err != nil {
		return err
	}
	return nil

}

func (s *CollectionService) GetImages(ctx context.Context, set *m.Collection, pag g.PaginationParams) ([]m.Image, *g.PaginationMeta, error) {

	return nil, nil, nil

}
