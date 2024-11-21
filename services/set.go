package services

import (
	"context"
	"github.com/google/uuid"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
)

type SetService struct {
	SetRepo         r.SetRepo
	ImageRepo       r.ImageRepo
	ImageService    ImageService
	MaxPageSize     int
	DefaultPageSize int
}

func (s *SetService) Create(ctx context.Context, set *m.Set) (*m.Set, error) {
	if err := set.Validate(); err != nil {
		return nil, err
	}

	set.Id = uuid.New()
	set, err := s.SetRepo.Create(ctx, set)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (s *SetService) GetOne(ctx context.Context, id string) (*m.Set, error) {
	set, err := s.SetRepo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (s *SetService) AppendImageToSet(ctx context.Context, image *m.Image, set *m.Set) error {
	err := s.SetRepo.AssignImageToSet(ctx, image, set)
	if err != nil {
		return err
	}
	return nil

}

func (s *SetService) GetImages(ctx context.Context, set *m.Set, pag g.PaginationParams) ([]m.Image, *g.PaginationMeta, error) {

	return nil, nil, nil

}
