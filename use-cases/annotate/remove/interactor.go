package remove

import (
	"fmt"

	"context"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	sauth "github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	repo Repo
	auth auth.Auth
}

type Option func(*Interactor)

func WithAuth(a auth.Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(repo Repo, opts ...Option) Interactor {
	i := &Interactor{repo: repo,
		auth: sauth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "removing annotation"
	id, err := a.NewAnnotationIdFromString(r.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	group, err := i.repo.GroupOfAnnotation(*id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if group != nil {

		if err := i.auth.AnnotateGroup(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	if err := i.repo.RemoveAnnotation(*id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessDeleteAnnotation(Response{Id: *id})

}
