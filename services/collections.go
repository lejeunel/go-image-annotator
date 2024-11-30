package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
)

type CollectionService struct {
	ImageService    ImageService
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
		return err
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

func (s *CollectionService) Clone(ctx context.Context, collection *m.Collection, newCollection *m.Collection) error {
	if err := g.CheckAuthorization(ctx, "admin"); err != nil {
		return err
	}

	s.Create(ctx, newCollection)

	filters := &g.ImageFilterArgs{}
	page := 1
	hasNextPage := true
	for hasNextPage {
		images, pageMeta, err := s.ImageService.GetPage(ctx, g.PaginationParams{Page: int64(page), PageSize: 1},
			filters, false)
		if err != nil {
			return err
		}
		s.CollectionRepo.AssignImageToCollection(ctx, &images[0], newCollection)
		if page == pageMeta.TotalPages {
			break
		}
		page += 1

	}

	return nil

}

func (s *CollectionService) Merge(ctx context.Context, source *m.Collection, destination *m.Collection) error {

	filters := &g.ImageFilterArgs{CollectionId: source.Id.String()}
	page := 1
	hasNextPage := true
	for hasNextPage {
		images, pageMeta, err := s.ImageService.GetPage(ctx, g.PaginationParams{Page: int64(page), PageSize: 1},
			filters, false)
		if err != nil {
			return err
		}
		imageFoundInDestination, err := s.CollectionRepo.ImageIsInCollection(ctx, &images[0], destination)
		if err != nil {
			return err
		}
		if !imageFoundInDestination {
			fmt.Println("merging new image")
			s.CollectionRepo.AssignImageToCollection(ctx, &images[0], destination)
		}
		if page == pageMeta.TotalPages {
			break
		}
		page += 1

	}

	return nil
}
