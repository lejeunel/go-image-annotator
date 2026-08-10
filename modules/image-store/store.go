package image_store

import (
	"fmt"
	"strings"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	fs "github.com/lejeunel/go-image-annotator/modules/file-store"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Transactor interface {
	RunInTx(fn func(Repos) error) error
}

type Repos struct {
	ImageRepo
	CollectionRepo
	AnnotationRepo
	MetaRepo
}

type ImageStore struct {
	Repos
	Transactor
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

	labels, err := s.AnnotationRepo.FindImageLabels(base.ImageId, collection.Name)
	if err != nil {
		return nil, fmt.Errorf("fetching labels: %w", err)
	}

	boxes, err := s.AnnotationRepo.FindBoundingBoxes(base.ImageId, collection.Name)
	if err != nil {
		return nil, fmt.Errorf("fetching bounding boxes: %w", err)
	}

	polygons, err := s.AnnotationRepo.FindPolygons(base.ImageId, collection.Name)
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
func (s ImageStore) Delete(id im.ImageId, collection clc.CollectionName) error {
	errCtx := fmt.Errorf("deleting image")
	if err := s.Transactor.RunInTx(func(tx Repos) error {
		if err := tx.AnnotationRepo.RemoveAllAnnotations(id, collection); err != nil {
			return fmt.Errorf("%w: %w", errCtx, err)
		}
		if err := tx.MetaRepo.DeleteAll(collection, id); err != nil {
			return fmt.Errorf("%w: %w", errCtx, err)
		}
		if err := tx.ImageRepo.RemoveImageFromCollection(id, collection); err != nil {
			return fmt.Errorf("%w: %w", errCtx, err)
		}
		return nil
	}); err != nil {
		return err
	}

	isUsed, err := s.ImageRepo.IsUsed(id)
	if err != nil {
		return fmt.Errorf(
			"%w: checking whether image %v is used in another collection: %w",
			errCtx,
			id,
			err,
		)
	}
	if !*isUsed {
		if err := s.DeleteAsset(id); err != nil {
			return fmt.Errorf("%w: %w", errCtx, err)
		}
	}
	return nil
}
func (s ImageStore) Copy(src clc.CollectionName, id im.ImageId, dst clc.CollectionName, deep bool) error {
	errCtx := fmt.Errorf("copying image %v from collection %v to collection %v", id, src, dst)
	image, err := s.Find(im.BaseImage{ImageId: id, Collection: src})
	if err != nil {
		return fmt.Errorf("%w: finding source image: %w", errCtx, err)
	}

	if err := s.Transactor.RunInTx(func(tx Repos) error {
		if err := tx.ImageRepo.AddToCollection(image.Id, dst); err != nil {
			return fmt.Errorf("%w: adding image to collection: %w", errCtx, err)
		}
		if deep {
			for _, label := range image.Labels {
				label.Id = a.NewAnnotationId()
				if err := tx.AnnotationRepo.AddImageLabel(
					image.Id,
					dst,
					label,
					label.Author,
					label.Time,
				); err != nil {
					return fmt.Errorf("%w: adding image label: %w", errCtx, err)
				}
			}

			for _, box := range image.BoundingBoxes {
				box.Id = a.NewAnnotationId()
				if err := tx.AnnotationRepo.AddBoundingBox(
					image.Id,
					dst,
					box,
					box.Author,
					box.Time,
				); err != nil {
					return fmt.Errorf("%w: adding bounding boxes: %w", errCtx, err)
				}
			}
			for _, poly := range image.Polygons {
				poly.Id = a.NewAnnotationId()
				if err := tx.AnnotationRepo.AddPolygon(
					image.Id,
					dst,
					poly,
					poly.Author,
					poly.Time,
				); err != nil {
					return fmt.Errorf("%w: adding polygons: %w", errCtx, err)
				}
			}

			for _, m := range image.Meta {
				if err := tx.MetaRepo.Add(dst, image.Id, m.Key, m.Value); err != nil {
					return fmt.Errorf("%w: adding meta-data with key %v: %w", errCtx, m.Key, err)
				}
			}
		}
		return nil

	}); err != nil {
		return err

	}
	return nil

}

func New(r Repos, t Transactor, f fs.FileStore) ImageStore {
	return ImageStore{r, t, f}
}
