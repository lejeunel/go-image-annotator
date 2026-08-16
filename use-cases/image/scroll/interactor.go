package scroll

import (
	"context"
	"errors"
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type FilterValidator interface {
	Validate(im.FilterStr) error
}

type OrderingValidator interface {
	Validate(im.OrderStr) error
}

type Interactor struct {
	ImageRepo
	FilterValidator
	OrderingValidator
}

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

func New(
	ir ImageRepo,
	fv FilterValidator,
	ov OrderingValidator,
) Interactor {
	return Interactor{ir, fv, ov}
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "scrolling images"

	id, err := im.NewImageIdFromString(r.CurrentImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: parsing image id %v: %w", errCtx, r.CurrentImageId, err))
		return
	}

	exists, err := i.ImageRepo.ImageExists(id)
	if err != nil {
		out.Error(fmt.Errorf("%v: checking existing of image %v: %w", errCtx, id, err))
		return
	}
	if !exists {
		out.Error(fmt.Errorf("%v: checking existing of image %v: %w", errCtx, id, e.ErrNotFound))
		return

	}
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

	next, errNext := i.ImageRepo.GetAdjacent(id, r.FilterStr, r.OrderStr, im.ScrollNext)
	prev, errPrev := i.ImageRepo.GetAdjacent(id, r.FilterStr, r.OrderStr, im.ScrollPrevious)
	err = errors.Join(errNext, errPrev)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessScroll(Response{Next: next, Prev: prev})
}
