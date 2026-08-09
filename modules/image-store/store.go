package image_store

import (
	"fmt"
	"strings"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type ImageStore struct {
	ImageRepo
	CollectionRepo
	AnnotationRepo
	MetaRepo
	fs.FileStore
}

func (s ImageStore) Find(base im.BaseImage) (*im.Image, error) {
	collection, err := s.CollectionRepo.Find(base.Collection)
	if err != nil {
		return nil, fmt.Errorf("fetching collection by name (%v): %w", base.Collection, err)
	}

	ok, err := s.ImageRepo.ImageExistsInCollection(base.ImageId, collection.Name)
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

	meta, err := s.MetaRepo.List(base.Collection, base.ImageId)
	if err != nil {
		return nil, fmt.Errorf("fetching image meta-data: %w", err)
	}

	reader, err := s.FileStore.Get(
		fmt.Sprintf("%v.%v", base.ImageId, strings.Split(specs.MIMEType, "/")[1]),
	)
	if err != nil {
		return nil, fmt.Errorf("fetching raw data: %w", err)
	}
	return &im.Image{
		Id:         base.ImageId,
		Collection: *collection, Labels: labels,
		BoundingBoxes: boxes,
		Polygons:      polygons,
		Specs:         *specs,
		Meta:          meta,
		Reader:        reader,
	}, nil
}

func (s ImageStore) DeleteAsset(id im.ImageId) error {
	specs, err := s.ImageRepo.GetSpecs(id)
	if err != nil {
		return fmt.Errorf("fetching image specs")
	}
	return s.FileStore.Delete(fmt.Sprintf("%v.%v", id, strings.Split(specs.MIMEType, "/")[1]))
}

func New(i ImageRepo, c CollectionRepo, a AnnotationRepo, m MetaRepo, f fs.FileStore) ImageStore {
	return ImageStore{i, c, a, m, f}
}
