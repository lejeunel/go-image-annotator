package add_bbox

import (
	"context"
	"fmt"

	st "github.com/lejeunel/go-image-annotator/app/image-store"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	sauth "github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"log/slog"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	imageStore st.Interface
	repo       Repo
	logger     *slog.Logger
	auth       auth.Auth
}

func NewInteractor(imageStore st.Interface, repo Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: repo, imageStore: imageStore, logger: logging.NewNoOpLogger(),
		auth: sauth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

type Option func(*Interactor)

func WithAuth(a auth.Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	image, err := i.findImage(r.ImageId, r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if image.Collection.Group != nil {
		if err := i.auth.AnnotateGroup(ctx, image.Collection.Group.Name); err != nil {
			i.handleError(err, out)
			return
		}
	}

	label, err := i.findLabel(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	box := a.NewBoundingBox(a.NewAnnotationId(), r.Xc, r.Yc, r.Width, r.Height, *label)
	if err := i.validateBox(image, box); err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.addBox(image, box); err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessAddBox(box)

}
func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "adding bounding box"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
func (i *Interactor) addBox(image *im.Image, box a.BoundingBox) error {
	if err := i.repo.AddBoundingBox(image.Id, image.Collection.Id, box); err != nil {
		return err
	}
	return nil
}

func (i *Interactor) validateBox(image *im.Image, box a.BoundingBox) error {
	if err := image.AddBoundingBox(box); err != nil {
		return err
	}
	return nil
}

func (i *Interactor) findLabel(name string) (*lbl.Label, error) {
	label, err := i.repo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil
}

func (i *Interactor) findImage(imageId string, collectionName string) (*im.Image, error) {
	image, err := i.imageStore.Find(im.BaseImage{ImageId: imageId, Collection: collectionName})
	if err != nil {
		return nil, err
	}
	return image, nil
}
