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

func (s *CollectionService) appendAnnotationsToImage(ctx context.Context, image *m.Image, collection *m.Collection) error {
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
func (s *CollectionService) AssignImageToCollection(ctx context.Context, image *m.Image, collection *m.Collection) error {
	return s.CollectionRepo.AssignImageToCollection(ctx, image, collection)
}

func (s *CollectionService) Clone(ctx context.Context, collection *m.Collection, newCollection *m.Collection, deep bool) error {
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
		if err := s.AssignImageToCollection(ctx, image, newCollection); err != nil {
			return err
		}

		if deep {
			if err := s.appendAnnotationsToImage(ctx, image, newCollection); err != nil {
				return err
			}

		}

		if page == pageMeta.TotalPages {
			break
		}
		page += 1

	}

	return nil

}

func (s *CollectionService) Merge(ctx context.Context, source *m.Collection, destination *m.Collection, deep bool) error {

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
			if err := s.CollectionRepo.AssignImageToCollection(ctx, &images[0], destination); err != nil {
				return err
			}
		}
		if deep {
			if err := s.appendAnnotationsToImage(ctx, image, destination); err != nil {
				return err
			}

		}

		if page == pageMeta.TotalPages {
			break
		}
		page += 1

	}

	return nil
}

func (s *CollectionService) RemoveImage(ctx context.Context, image *m.Image, collection *m.Collection) error {
	return s.CollectionRepo.RemoveImage(ctx, image, collection)
}

func (s *CollectionService) GetOrdinal(ctx context.Context, index int) (*m.Collection, error) {
	collections, _, err := s.GetPage(ctx,
		g.PaginationParams{Page: int64(index), PageSize: 1})
	if err != nil {
		return nil, err
	}
	return &collections[0], nil

}

func (s *CollectionService) GetPage(
	ctx context.Context,
	pagination g.PaginationParams) ([]m.Collection, *g.PaginationMeta, error) {

	p := s.CollectionRepo.Paginate(pagination.PageSize)
	p.SetPage(int(pagination.Page))

	var collections []m.Collection
	err := p.Results(&collections)

	if err != nil {
		return nil, nil, err
	}

	paginationMeta := g.NewPaginationMeta(p)
	return collections, &paginationMeta, nil

}
