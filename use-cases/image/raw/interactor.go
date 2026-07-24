package raw

import (
	"fmt"
	"io"
	"strings"

	im "github.com/lejeunel/go-image-annotator/entities/image"
)

type Interface interface {
	Execute(Request, OutputPort)
}

type FileGetter interface {
	Get(string) (io.Reader, error)
}

type Repo interface {
	GetSpecs(im.ImageId) (*im.ImageSpecs, error)
}

type Interactor struct {
	fileGetter FileGetter
	repo       Repo
}

func New(fileGetter FileGetter, repo Repo) Interactor {
	return Interactor{fileGetter: fileGetter, repo: repo}
}

func (i Interactor) Execute(id string, out OutputPort) {
	errCtx := "reading raw image data"
	imageId, err := im.NewImageIdFromString(id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	specs, err := i.repo.GetSpecs(imageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching image specifications: %w", errCtx, err))
		return
	}
	reader, err := i.fileGetter.Get(fmt.Sprintf("%v.%v", imageId.String(), strings.Split(specs.MIMEType, "/")[1]))
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching raw-data: %w", errCtx, err))
		return
	}

	out.SuccessReadRawImage(Response{Reader: reader, ImageSpecs: *specs})
}
