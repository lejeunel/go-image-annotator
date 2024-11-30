package services

import (
	"context"
	"github.com/google/uuid"
	g "go-image-annotator/generic"
	m "go-image-annotator/models"
	r "go-image-annotator/repositories"
)

type CollectionService struct {
	ImageService      ImageService
	AnnotationService AnnotationService
	CollectionRepo    r.CollectionRepo
	ImageRepo         r.ImageRepo
	MaxPageSize       int
	DefaultPageSize   int
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

func (s *CollectionService) appendImage(ctx context.Context, image *m.Image, collection *m.Collection, deep bool) error {
	if image.Annotations != nil {
		for _, a := range image.Annotations {
			if err := s.AnnotationService.ApplyLabelToImage(ctx, a.Label, image, collection); err != nil {
				return err
			}

		}
	}
	if image.BoundingBoxes != nil {
		for _, bbox := range image.BoundingBoxes {

			if err := s.AnnotationService.ApplyBoundingBoxToImage(ctx, bbox, image, collection); err != nil {
				return err
			}

		}
	}
	return nil

}

func (s *CollectionService) Clone(ctx context.Context, collection *m.Collection, newCollection *m.Collection) error {
	if err := g.CheckAuthorization(ctx, "admin"); err != nil {
		return err
	}

	s.Create(ctx, newCollection)

	page := 1
	hasNextPage := true
	for hasNextPage {
		images, pageMeta, err := s.ImageService.GetPage(ctx, collection.Id.String(),
			g.PaginationParams{Page: int64(page), PageSize: 1},
			false)
		image := &images[0]
		if err != nil {
			return err
		}
		err = s.CollectionRepo.AssignImageToCollection(ctx, image, newCollection)
		if err != nil {
			return err
		}

		if err := s.appendImage(ctx, image, newCollection, true); err != nil {
			return err
		}

		if page == pageMeta.TotalPages {
			break
		}
		page += 1

	}

	return nil

}

func (s *CollectionService) Merge(ctx context.Context, source *m.Collection, destination *m.Collection) error {

	page := 1
	hasNextPage := true
	for hasNextPage {
		images, pageMeta, err := s.ImageService.GetPage(ctx, source.Id.String(), g.PaginationParams{Page: int64(page), PageSize: 1},
			false)
		image := &images[0]
		if err != nil {
			return err
		}
		imageFoundInDestination, err := s.CollectionRepo.ImageIsInCollection(ctx, &images[0], destination)
		if err != nil {
			return err
		}
		if !imageFoundInDestination {
			s.CollectionRepo.AssignImageToCollection(ctx, &images[0], destination)
		}
		if err := s.appendImage(ctx, image, destination, true); err != nil {
			return err
		}

		if page == pageMeta.TotalPages {
			break
		}
		page += 1

	}

	return nil
}
