package services

import (
	"context"
	"github.com/google/uuid"
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

func (s *CollectionService) Create(ctx context.Context, collection *m.Collection) error {
	if err := collection.Validate(); err != nil {
		return err
	}

	collection.Id = uuid.New()
	collection, err := s.CollectionRepo.Create(ctx, collection)
	if err != nil {
		return nil
	}
	return nil
}

func (s *CollectionService) GetOne(ctx context.Context, id string) (*m.Collection, error) {
	set, err := s.CollectionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return set, nil
}
