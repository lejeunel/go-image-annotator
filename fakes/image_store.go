package fake

import (
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ImageStore struct {
	ErrOnFind               error
	ErrOnDelete             error
	ErrOnPaginateCollection error
	ErrOnSlice              error
	Return                  *im.Image
	DeletedAssetId          *im.ImageId
	DeletedId               *im.ImageId
	DeletedBatch            bool
	CopiedToCollection      string
	ReturnPaginated         []im.Image
	ReturnCount             int64
	ReturnPagination        pag.Pagination
	GotPagination           pag.PaginationParams
}

func (s *ImageStore) Find(baseImage im.BaseImage) (*im.Image, error) {
	if s.ErrOnFind != nil {
		return nil, s.ErrOnFind
	}
	if s.Return != nil {
		return s.Return, nil
	}
	return &im.Image{}, nil
}

func (s *ImageStore) PaginateCollection(
	name clc.CollectionName,
	pag pag.PaginationParams,
) ([]im.Image, *int64, error) {
	if s.ErrOnPaginateCollection != nil {
		return nil, nil, s.ErrOnPaginateCollection
	}
	s.GotPagination = pag
	return s.ReturnPaginated, &s.ReturnCount, nil
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

func (s *ImageStore) Copy(
	src clc.CollectionName,
	id im.ImageId,
	dst clc.CollectionName,
	deep bool,
) error {
	s.CopiedToCollection = dst
	return nil
}
