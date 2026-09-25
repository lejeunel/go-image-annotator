package list

import (
	"fmt"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type ImageStore interface {
	PaginateCollection(
		clc.CollectionName,
		pagination.PaginationParams,
	) ([]im.Image, *pagination.Pagination, error)
}

type Interactor struct {
	ImageStore
	DefaultPageSize int
	MaxPageSize     int
}

func New(
	s ImageStore,
	dps int,
	mps int,
) Interactor {
	return Interactor{s, dps, mps}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := fmt.Errorf("listing images")

	r.PaginationParams.Sanitize(i.DefaultPageSize, i.MaxPageSize)

	images, pagination, err := i.ImageStore.PaginateCollection(r.CollectionName, r.PaginationParams)
	if err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, err))
		return
	}
	response := Response{
		Images:     images,
		Pagination: *pagination,
	}

	out.SuccessPaginateImages(response)
}
