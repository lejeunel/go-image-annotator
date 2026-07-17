package image_store

import (
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type ImageStore struct {
	ImageRepo
	CollectionRepo
	AnnotationRepo
	fileStore fs.Interface
}

func (s ImageStore) Find(base im.BaseImage) (*im.Image, error) {
	collection, err := s.CollectionRepo.FindCollectionByName(base.Collection)
	if err != nil {
		return nil, fmt.Errorf("fetching collection by name (%v): %w", base.Collection, err)
	}

	ok, err := s.ImageRepo.ImageExistsInCollection(base.ImageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("checking whether image %v exists in collection %v: %w",
			base.ImageId, base.Collection, err)
	}
	if !ok {
		return nil, fmt.Errorf("checking whether image %v exists in collection %v: %w",
			base.ImageId, base.Collection, e.ErrNotFound)

	}

	labels, err := s.AnnotationRepo.FindImageLabels(base.ImageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("fetching labels: %w", err)
	}

	boxes, err := s.AnnotationRepo.FindBoundingBoxes(base.ImageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("fetching bounding boxes: %w", err)
	}

	polygons, err := s.AnnotationRepo.FindPolygons(base.ImageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("fetching polygons: %w", err)
	}

	specs, err := s.ImageRepo.GetSpecs(base.ImageId)
	if err != nil {
		return nil, fmt.Errorf("fetching image specs: %w", err)
	}

	reader, err := s.fileStore.Get(base.ImageId)
	if err != nil {
		return nil, fmt.Errorf("fetching raw data: %w", err)
	}
	return &im.Image{
		Id:         base.ImageId,
		Collection: *collection, Labels: labels,
		BoundingBoxes: boxes,
		Polygons:      polygons,
		Specs:         *specs,
		Reader:        reader}, nil

}

func New(i ImageRepo, c CollectionRepo, a AnnotationRepo, f fs.Interface) ImageStore {
	return ImageStore{i, c, a, f}
}
