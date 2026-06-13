package image_store

import (
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type ImageStore struct {
	repo      Repo
	fileStore fs.Interface
}

func (s ImageStore) Find(base im.BaseImage) (*im.Image, error) {
	imageId, err := im.NewImageIdFromString(base.ImageId)
	if err != nil {
		return nil, fmt.Errorf("fetching collection by name (%v): %w", base.Collection, err)
	}

	collection, err := s.repo.FindCollectionByName(base.Collection)
	if err != nil {
		return nil, fmt.Errorf("fetching collection by name (%v): %w", base.Collection, err)
	}

	ok, err := s.repo.ImageExistsInCollection(imageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("checking whether image %v exists in collection %v: %w",
			imageId, base.Collection, err)
	}
	if !ok {
		return nil, fmt.Errorf("checking whether image %v exists in collection %v: %w",
			imageId, base.Collection, e.ErrNotFound)

	}

	labels, err := s.repo.FindImageLabels(imageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("fetching labels: %w", err)
	}

	boxes, err := s.repo.FindBoundingBoxes(imageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("fetching bounding boxes: %w", err)
	}

	specs, err := s.repo.GetSpecs(imageId)
	if err != nil {
		return nil, fmt.Errorf("fetching image specs: %w", err)
	}

	reader, err := s.fileStore.Get(imageId)
	if err != nil {
		return nil, fmt.Errorf("fetching raw data: %w", err)
	}
	return &im.Image{Id: imageId, Collection: *collection, Labels: labels,
		BoundingBoxes: boxes,
		Specs:         *specs,
		Reader:        reader}, nil

}

func New(repo Repo, fileStore fs.Interface) ImageStore {
	return ImageStore{repo: repo, fileStore: fileStore}
}
