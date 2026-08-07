package list

import (
	"context"
	"fmt"

	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	MetaDataRepo
}

func New(m MetaDataRepo) Interactor {
	return Interactor{MetaDataRepo: m}
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("listing metadata for image %v and collection %v", r.ImageId, r.Collection)

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: parsing image id %v: %w", errCtx, imageId, err))
		return
	}
	meta, err := i.MetaDataRepo.List(r.Collection, imageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w",
			errCtx, err))
		return
	}

	out.SuccessListMetadata(meta)
}
