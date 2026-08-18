package fake

import (
	"iter"
	"slices"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ImageRepo struct {
	RemovedImageIdFromCollection *im.ImageId
	ErrOnRemoveImage             error
	ErrOnImageExistsInCollection error
	ErrOnImageExists             error
	ErrOnGetSpecs                error
	ErrOnAddToCollection         error
	ErrOnSlice                   error
	ErrOnAddImage                error
	ErrOnDeleteImage             error
	ErrOnFindHash                error
	ErrOnCount                   error
	ErrOnIterate                 error
	ErrOnIsUsed                  error
	ErrOnGetAdjacent             error
	AddedImageId                 im.ImageId
	AddedIntoCollection          *clc.CollectionName
	ImageIsInCollection          bool
	IsUsed_                      bool
	GotFilters                   im.FilterStr
	GotPagination                pa.PaginationParams
	GotOrdering                  im.OrderStr
	GotHash                      []byte
	GotSpecs                     im.Specs
	ReturnSpecs                  *im.Specs
	NumDeletedImages             int
	HashAlreadyExists            bool
	Count_                       int64
	IterateBaseImages            []im.BaseImage
	NextImage                    im.BaseImage
	PreviousImage                im.BaseImage
	ImageMissing                 bool
}

func (r *ImageRepo) RemoveImageFromCollection(
	imageId im.ImageId,
	collection clc.CollectionName,
) error {
	if r.ErrOnRemoveImage != nil {
		return r.ErrOnRemoveImage
	}
	r.RemovedImageIdFromCollection = &imageId
	return nil
}

func (r *ImageRepo) ImageExists(imageId im.ImageId) (bool, error) {
	if r.ErrOnImageExists != nil {
		return false, r.ErrOnImageExists
	}
	if r.ImageMissing {
		return false, nil
	}
	return true, nil
}

func (r *ImageRepo) ImageExistsInCollection(
	imageId im.ImageId,
	collection clc.CollectionName,
) (bool, error) {
	if r.ErrOnImageExistsInCollection != nil {
		return false, r.ErrOnImageExistsInCollection
	}
	if r.ImageIsInCollection {
		return true, nil
	}
	return false, nil
}

func (r *ImageRepo) AddToCollection(imageId im.ImageId, collection clc.CollectionName) error {
	if r.ErrOnAddToCollection != nil {
		return r.ErrOnAddToCollection
	}
	r.AddedImageId = imageId
	r.AddedIntoCollection = &collection
	return nil
}

func (r *ImageRepo) Slice(
	f im.FilterStr,
	p pa.PaginationParams,
	o im.OrderStr,
) ([]im.BaseImage, error) {
	r.GotFilters = f
	r.GotPagination = p
	r.GotOrdering = o

	if r.ErrOnSlice != nil {
		return nil, r.ErrOnSlice
	}

	result := []im.BaseImage{}
	collectionName := "a-collection"
	for range p.PageSize {
		result = append(result,
			im.BaseImage{
				Collection: collectionName,
				ImageId:    im.NewImageId(),
			})
	}

	return result, nil
}

func (r *ImageRepo) AddImage(imageId im.ImageId, hash []byte, specs im.Specs) error {
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

func (r *ImageRepo) Count(f im.FilterStr) (*int64, error) {
	if r.ErrOnCount != nil {
		return nil, r.ErrOnCount
	}
	return &r.Count_, nil
}

func (r ImageRepo) GetSpecs(im.ImageId) (*im.Specs, error) {
	if r.ErrOnGetSpecs != nil {
		return nil, r.ErrOnGetSpecs
	}
	return r.ReturnSpecs, nil
}

func (r ImageRepo) Iterate(f im.FilterStr, pageSize int) iter.Seq2[im.BaseImage, error] {
	return func(yield func(im.BaseImage, error) bool) {
		for img := range slices.Values(r.IterateBaseImages) {
			if !yield(img, nil) {
				return
			}
		}
	}
}

func (r *ImageRepo) IsUsed(id im.ImageId) (*bool, error) {
	if r.ErrOnIsUsed != nil {
		return nil, r.ErrOnIsUsed
	}
	return &r.IsUsed_, nil
}

func (r *ImageRepo) GetAdjacent(
	id im.ImageId,
	c string,
	f im.FilterStr,
	o im.OrderStr,
	d im.ScrollingDirection,
) (*im.BaseImage, error) {
	if r.ErrOnGetAdjacent != nil {
		return nil, r.ErrOnGetAdjacent
	}
	if d == im.ScrollNext {
		return &r.NextImage, nil
	}
	return &r.PreviousImage, nil

}
