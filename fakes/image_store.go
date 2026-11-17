package fake

import (
	"iter"
	"slices"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ImageStore struct {
	ErrOnFind                 error
	ErrOnDelete               error
	ErrOnPaginateCollection   error
	ErrOnSlice                error
	Return                    *im.Image
	DeletedAssetId            *im.ImageId
	DeletedId                 *im.ImageId
	DeletedBatch              bool
	DeepCopiedToCollection    string
	ShallowCopiedToCollection string
	ReturnPaginated           []im.Image
	ReturnCount               int64
	ReturnPagination          pag.Pagination
	GotPagination             pag.PaginationParams
	IterateImageIds           []im.ImageId
	ExistingImageFileIds      []im.ImageId
	ErrOnIterateImageIDs      error
	ErrOnExists               error
	// Existing, when non-nil, makes Find resolve images per collection:
	// only listed (image, collection) pairs are found, and copies add to it.
	Existing []im.BaseImage
}

func (s *ImageStore) Find(baseImage im.BaseImage) (*im.Image, error) {
	if s.ErrOnFind != nil {
		return nil, s.ErrOnFind
	}
	if s.Existing != nil {
		if !slices.Contains(s.Existing, baseImage) {
			return nil, e.ErrNotFound
		}
		return &im.Image{
			Id:         baseImage.ImageId,
			Collection: clc.Collection{Name: baseImage.Collection},
		}, nil
	}
	if s.Return != nil {
		return s.Return, nil
	}
	return &im.Image{}, nil
}

func (s *ImageStore) PaginateCollection(
	name clc.CollectionName,
	pag pag.PaginationParams,
) ([]im.Image, *pag.Pagination, error) {
	if s.ErrOnPaginateCollection != nil {
		return nil, nil, s.ErrOnPaginateCollection
	}
	s.GotPagination = pag
	return s.ReturnPaginated, &s.ReturnPagination, nil
}

func (s *ImageStore) Slice(
	f im.FilterStr,
	o im.OrderStr,
	pag pag.PaginationParams,
) ([]im.Image, *pag.Pagination, error) {
	if s.ErrOnSlice != nil {
		return nil, nil, s.ErrOnSlice
	}
	s.GotPagination = pag
	return s.ReturnPaginated, &s.ReturnPagination, nil
}

func (s *ImageStore) DeleteAsset(id im.ImageId) error {
	s.DeletedAssetId = &id
	return nil
}

func (s *ImageStore) Delete(id im.ImageId, collection clc.CollectionName) error {
	if s.ErrOnDelete != nil {
		return s.ErrOnDelete
	}
	s.DeletedId = &id
	return nil
}

func (s *ImageStore) DeleteBatch([]im.ImageId, clc.CollectionName) error {
	s.DeletedBatch = true
	return nil
}

func (s *ImageStore) DeepCopy(
	src clc.CollectionName,
	id im.ImageId,
	dst clc.CollectionName,
) error {
	s.DeepCopiedToCollection = dst
	s.recordCopy(id, dst)
	return nil
}

func (s *ImageStore) ShallowCopy(
	id im.ImageId,
	dst clc.CollectionName,
) error {
	s.ShallowCopiedToCollection = dst
	s.recordCopy(id, dst)
	return nil
}

func (s ImageStore) IterateImageIDs(pageSize int) iter.Seq2[im.ImageId, error] {
	return func(yield func(im.ImageId, error) bool) {
		if s.ErrOnIterateImageIDs != nil {
			yield(im.ImageId{}, s.ErrOnIterateImageIDs)
			return
		}
		for id := range slices.Values(s.IterateImageIds) {
			if !yield(id, nil) {
				return
			}
		}
	}
}

func (s ImageStore) Exists(id im.ImageId) (*bool, error) {
	if s.ErrOnExists != nil {
		return nil, s.ErrOnExists
	}
	result := false
	if slices.Contains(s.ExistingImageFileIds, id) {
		result = true
		return &result, nil
	}
	return &result, nil
}

func (s *ImageStore) recordCopy(id im.ImageId, dst clc.CollectionName) {
	if s.Existing != nil {
		s.Existing = append(s.Existing, im.BaseImage{ImageId: id, Collection: dst})
	}
}
