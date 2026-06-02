package modify_bbox

import (
	"context"
	"fmt"

	"log/slog"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	sauth "github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	repo   Repo
	logger *slog.Logger
	auth   auth.Auth
}

type Option func(*Interactor)

func WithAuth(a auth.Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func NewInteractor(repo Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: repo, logger: logging.NewNoOpLogger(),
		auth: sauth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	annotationId, err := a.NewAnnotationIdFromString(r.AnnotationId)
	if err != nil {
		i.handleError(err, out)
		return
	}
	group, err := i.repo.GroupOfAnnotation(*annotationId)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.auth.AnnotateGroup(ctx, *group); err != nil {
		i.handleError(err, out)
		return
	}
	label, err := i.findLabel(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}
	u, err := i.validate(r.Xc, r.Yc, r.Width, r.Height, *label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.update(*annotationId, *u); err != nil {
		i.handleError(err, out)
		return
	}
	out.SuccessUpdateBox(Response{})

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "updating bounding box properties"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func (i *Interactor) update(id a.AnnotationId, u a.BoundingBoxUpdatables) error {
	if err := i.repo.UpdateBoundingBox(id, u); err != nil {
		return err
	}
	return nil

}

func (i *Interactor) validate(xc float32, yc float32, width float32,
	height float32, label lbl.Label) (*a.BoundingBoxUpdatables, error) {

	if err := a.ValidateBoundingBox(xc, yc, width, height); err != nil {
		return nil, err
	}
	return &a.BoundingBoxUpdatables{LabelId: label.Id, Xc: xc, Yc: yc, Width: width, Height: height}, nil

}

func (i *Interactor) findLabel(name string) (*lbl.Label, error) {

	label, err := i.repo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil

}
