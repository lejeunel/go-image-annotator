package list

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
	Find(base im.BaseImage) (*im.Image, error)
}

type Interactor struct {
	Repo
	FilterValidator
	OrderingValidator
	ImageStore
	DefaultPageSize int
	MaxPageSize     int
}

func New(
	r Repo,
	fv FilterValidator,
	ov OrderingValidator,
	s ImageStore,
	dps int,
	mps int,
) Interactor {
	return Interactor{r, fv, ov, s, dps, mps}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := "listing images"

	r.PaginationParams.Sanitize(i.DefaultPageSize, i.MaxPageSize)

	if r.FilterStr != "" {
		if err := i.FilterValidator.Validate(r.FilterStr); err != nil {
			out.Error(fmt.Errorf("%v: validating query %v: %w", errCtx, r.FilterStr, err))
			return
		}
	}

	if r.OrderStr != "" {
		if err := i.OrderingValidator.Validate(r.OrderStr); err != nil {
			out.Error(fmt.Errorf("%v: validating ordering %v: %w", errCtx, r.OrderStr, err))
			return
		}
	}

	baseImages, err := i.Repo.Slice(r.FilterStr, r.PaginationParams, r.OrderStr)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	count, err := i.Repo.Count(r.FilterStr)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	imageResponses, err := i.buildResponse(baseImages)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	response := Response{
		Images: imageResponses,
		Pagination: pagination.Pagination{
			Page:         r.Page,
			PageSize:     r.PageSize,
			TotalRecords: *count,
			TotalPages:   (*count + int64(r.PageSize) - 1) / int64(r.PageSize),
		},
	}

	out.SuccessListImages(response)
}

func (i *Interactor) buildResponse(baseImages []im.BaseImage) ([]im.Image, error) {
	r := []im.Image{}
	for _, baseImage := range baseImages {
		image, err := i.ImageStore.Find(baseImage)
		if err != nil {
			return nil, err
		}
		r = append(r, *image)
	}
	return r, nil
}
