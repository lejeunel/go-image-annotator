package list

import (
	"fmt"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	ist "github.com/lejeunel/go-image-annotator/modules/image-store"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	repo  Repo
	store ist.Interface
}

func New(r Repo, s ist.Interface) Interactor {
	return Interactor{repo: r, store: s}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := "listing images"

	if r.PageSize == 0 {
		out.Error(fmt.Errorf("%v: got page size = 0", errCtx))
		return
	}

	baseImages, err := i.repo.List(r.FilteringParams, r.OrderingParams)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	count, err := i.repo.Count(im.CountingParams{Collection: r.FilteringParams.Collection})
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	imageResponses, err := i.buildResponse(baseImages)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	response := Response{Images: imageResponses,
		Pagination: pagination.Pagination{Page: r.Page, PageSize: r.PageSize, TotalRecords: *count, TotalPages: *count / int64(r.PageSize)}}

	out.Success(response)

}

func (i *Interactor) buildResponse(baseImages []im.BaseImage) ([]im.Image, error) {
	r := []im.Image{}
	for _, baseImage := range baseImages {
		image, err := i.store.Find(baseImage)
		if err != nil {
			return nil, err
		}
		r = append(r, *image)
	}
	return r, nil

}
