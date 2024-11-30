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

func (s *CollectionService) Get(ctx context.Context, id string) (*m.Collection, error) {
	set, err := s.CollectionRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (s *CollectionService) Delete(ctx context.Context, collection *m.Collection) error {
	if err := g.CheckAuthorization(ctx, "admin"); err != nil {
		return err
	}

	if err := s.ImageRepo.DeleteImagesInCollection(ctx, collection); err != nil {
		return err
	}
	return s.CollectionRepo.Delete(ctx, collection)
}
