package slice

import (
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type FilterValidator interface {
	Validate(im.FilterStr) error
}

type OrderingValidator interface {
	Validate(im.OrderStr) error
}

type ImageStore interface {
	Slice(
		im.FilterStr,
		im.OrderStr,
		pagination.PaginationParams,
	) ([]im.Image, *pagination.Pagination, error)
}

type Interactor struct {
	ImageStore
	FilterValidator
	OrderingValidator
	DefaultPageSize int
	MaxPageSize     int
}

func New(
	s ImageStore,
	fv FilterValidator,
	ov OrderingValidator,
	dps int,
	mps int,
) Interactor {
	return Interactor{s, fv, ov, dps, mps}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := "slicing images"

	r.PaginationParams.Sanitize(i.DefaultPageSize, i.MaxPageSize)

	if r.FilterStr != "" {
		if err := i.FilterValidator.Validate(r.FilterStr); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	if r.OrderStr != "" {
		if err := i.OrderingValidator.Validate(r.OrderStr); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	images, pagination, err := i.ImageStore.Slice(r.FilterStr, r.OrderStr, r.PaginationParams)
	if err != nil {
		out.Error(err)
		return
	}

	out.SuccessSliceImages(
		Response{
			Images:     images,
			Pagination: *pagination,
			FilterStr:  r.FilterStr,
			OrderStr:   r.OrderStr,
		},
	)
}
