package image_store

import (
	"fmt"

	fs "github.com/lejeunel/go-image-annotator-v2/app/file-store"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

type ImageStore struct {
	repo      Repo
	fileStore fs.Interface
}

func (s ImageStore) Find(base im.BaseImage) (*im.Image, error) {
	collection, err := s.repo.FindCollectionByName(base.Collection)
	if err != nil {
		return nil, fmt.Errorf("fetching collection by name (%v): %w", base.Collection, err)
	}

	ok, err := s.repo.ImageExistsInCollection(base.ImageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("checking whether image %v exists in collection %v: %w",
			base.ImageId.String(), base.Collection, err)
	}
	if !ok {
		return nil, fmt.Errorf("checking whether image %v exists in collection %v: %w",
			base.ImageId.String(), base.Collection, e.ErrNotFound)

	}

	labels, err := s.repo.FindImageLabels(base.ImageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("fetching labels: %w", err)
	}

	boxes, err := s.repo.FindBoundingBoxes(base.ImageId, collection.Id)
	if err != nil {
		return nil, fmt.Errorf("fetching bounding boxes: %w", err)
	}

	specs, err := s.repo.GetSpecs(base.ImageId)
	if err != nil {
		return nil, fmt.Errorf("fetching image specs: %w", err)
	}

	reader, err := s.fileStore.Get(base.ImageId)
	if err != nil {
		return nil, fmt.Errorf("fetching raw data: %w", err)
	}
	return &im.Image{Id: base.ImageId, Collection: *collection, Labels: labels,
		BoundingBoxes: boxes,
		Specs:         *specs,
		Reader:        reader}, nil

}

func New(repo Repo, fileStore fs.Interface) *ImageStore {
	return &ImageStore{repo: repo, fileStore: fileStore}
}
