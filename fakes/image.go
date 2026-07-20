package fake

import (
	"iter"
	"slices"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ImageRepo struct {
	RemovedImageId               im.ImageId
	ErrOnRemoveImage             error
	ErrOnImageExistsInCollection error
	ErrOnImageExists             error
	ErrOnGetSpecs                error
	ErrOnAddToCollection         error
	ErrOnList                    error
	ErrOnAddImage                error
	ErrOnDeleteImage             error
	ErrOnFindHash                error
	ErrOnCount                   error
	ErrOnIterate                 error
	AddedImageId                 im.ImageId
	AddedIntoCollectionId        clc.CollectionId
	ImageIsInCollection          bool
	GotFilters                   im.Filtering
	GotPagination                pa.PaginationParams
	GotOrdering                  im.Ordering
	GotHash                      []byte
	GotSpecs                     im.ImageSpecs
	ReturnSpecs                  im.ImageSpecs
	NumDeletedImages             int
	HashAlreadyExists            bool
	Count_                       int64
	IterateBaseImages            []im.BaseImage
}

func (r *ImageRepo) RemoveImageFromCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	if r.ErrOnRemoveImage != nil {
		return r.ErrOnRemoveImage
	}
	r.RemovedImageId = imageId
	return nil
}

func (r *ImageRepo) ImageExists(imageId im.ImageId) (bool, error) {
	if r.ErrOnImageExists != nil {
		return false, r.ErrOnImageExists
	}
	return true, nil
}
func (r *ImageRepo) ImageExistsInCollection(imageId im.ImageId, collectionId clc.CollectionId) (bool, error) {
	if r.ErrOnImageExistsInCollection != nil {
		return false, r.ErrOnImageExistsInCollection
	}
	if r.ImageIsInCollection {
		return true, nil
	}
	return false, nil
}

func (r *ImageRepo) AddToCollection(imageId im.ImageId, collectionId clc.CollectionId) error {
	if r.ErrOnAddToCollection != nil {
		return r.ErrOnAddToCollection
	}
	r.AddedImageId = imageId
	r.AddedIntoCollectionId = collectionId
	return nil
}
func (r *ImageRepo) Slice(f im.Filtering, p pa.PaginationParams, o im.Ordering) ([]im.BaseImage, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList
	}

	r.GotFilters = f
	r.GotPagination = p
	r.GotOrdering = o

	result := []im.BaseImage{}
	collectionName := "a-collection"
	for range p.PageSize {
		result = append(result,
			im.BaseImage{
				Collection: collectionName,
				ImageId:    im.NewImageId()})
	}

	return result, nil

}
func (r *ImageRepo) AddImage(imageId im.ImageId, hash []byte, specs im.ImageSpecs) error {
	if r.ErrOnAddImage != nil {
		return r.ErrOnAddImage
	}
	r.GotHash = hash
	r.GotSpecs = specs
	return nil
}

func (r *ImageRepo) Delete(im.ImageId) error {
	if r.ErrOnDeleteImage != nil {
		return r.ErrOnDeleteImage
	}
	r.NumDeletedImages += 1
	return nil
}

func (r *ImageRepo) FindImageIdByHash(hash []byte) (*im.ImageId, error) {
	if r.ErrOnFindHash != nil {
		return nil, r.ErrOnFindHash
	}
	if r.HashAlreadyExists {
		existingId := im.NewImageId()
		return &existingId, nil
	}
	return nil, nil
}
func (r *ImageRepo) Count(f im.Filtering) (*int64, error) {
	if r.ErrOnCount != nil {
		return nil, r.ErrOnCount
	}
	return &r.Count_, nil

}

func (r ImageRepo) GetSpecs(im.ImageId) (*im.ImageSpecs, error) {
	if r.ErrOnGetSpecs != nil {
		return nil, r.ErrOnGetSpecs
	}
	return &r.ReturnSpecs, nil
}

func (r ImageRepo) Iterate(f im.Filtering, pageSize int) iter.Seq2[im.BaseImage, error] {
	return func(yield func(im.BaseImage, error) bool) {
		for img := range slices.Values(r.IterateBaseImages) {
			if !yield(img, nil) {
				return
			}
		}
	}
}
