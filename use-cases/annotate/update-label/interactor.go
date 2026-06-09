package update_label

import (
	"context"
	"fmt"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	sauth "github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"log/slog"
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

func New(repo Repo, opts ...Option) Interactor {
	i := &Interactor{repo: repo, logger: logging.NewNoOpLogger(),
		auth: sauth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	label, err := i.repo.FindLabel(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	id, err := a.NewAnnotationIdFromString(r.AnnotationId)
	if err != nil {
		i.handleError(err, out)
		return
	}

	group, err := i.repo.GroupOfAnnotation(*id)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if err := i.auth.AnnotateGroup(ctx, *group); err != nil {
		i.handleError(err, out)
		return
	}

	err = i.repo.UpdateLabelOfAnnotation(*id, label.Id)
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessUpdateLabel(Response{})

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "updating bounding box properties"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}
