package list

import (
	"fmt"
	ist "github.com/lejeunel/go-image-annotator-v2/app/image-store"
	im "github.com/lejeunel/go-image-annotator-v2/entities/image"
	"github.com/lejeunel/go-image-annotator-v2/shared/logging"
	"github.com/lejeunel/go-image-annotator-v2/shared/pagination"
	"log/slog"
)

type Interactor struct {
	repo   Repo
	store  ist.Interface
	logger *slog.Logger
}

func NewInteractor(r Repo, s ist.Interface) *Interactor {
	return &Interactor{repo: r, store: s, logger: logging.NewNoOpLogger()}
}

func (i *Interactor) Execute(r Request, out OutputPort) {
	filteringParams := &ist.FilteringParams{
		Page:       r.Page,
		PageSize:   r.PageSize,
		Collection: r.CollectionName}

	baseImages, err := i.repo.List(*filteringParams)
	if err != nil {
		i.handleError(err, out)
		return
	}

	count, err := i.repo.Count(ist.CountingParams{Collection: filteringParams.Collection})
	if err != nil {
		i.handleError(err, out)
		return
	}

	imageResponses, err := i.buildResponse(*baseImages)
	if err != nil {
		i.handleError(err, out)
		return
	}

	response := Response{Images: *imageResponses,
		Pagination: pagination.Pagination{Page: r.Page, PageSize: r.PageSize, TotalRecords: *count, TotalPages: *count / int64(r.PageSize)}}

	out.Success(response)

}

func (i *Interactor) buildResponse(baseImages []im.BaseImage) (*[]im.Image, error) {
	r := []im.Image{}
	for _, baseImage := range baseImages {
		image, err := i.store.Find(baseImage)
		if err != nil {
			return nil, err
		}
		r = append(r, *image)
	}
	return &r, nil

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "listing images"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
