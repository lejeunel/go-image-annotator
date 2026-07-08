package raw

import (
	"fmt"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	"io"
)

type Interface interface {
	Execute(Request, OutputPort)
}

type FileGetter interface {
	Get(im.ImageId) (io.Reader, error)
}

type Interactor struct {
	fileGetter FileGetter
}

func New(fileGetter FileGetter) Interactor {
	return Interactor{fileGetter: fileGetter}
}

func (i Interactor) Execute(r Request, out OutputPort) {
	errCtx := "reading raw image data"
	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	reader, err := i.fileGetter.Get(imageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessReadRawImage(reader)
}
