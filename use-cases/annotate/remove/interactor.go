package remove

import (
	"context"
	"fmt"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	sauth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	Repo
	auth.Auth
}

type Option func(*Interactor)

func WithAuth(a auth.Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(repo Repo, opts ...Option) Interactor {
	i := &Interactor{
		Repo: repo,
		Auth: sauth.NewVoidAuth(),
	}
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

	group, err := i.Repo.GroupOfAnnotation(*id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if group != nil {
		if err := i.Auth.Annotate(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	if err := i.Repo.RemoveAnnotation(*id); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessDeleteAnnotation(Response{Id: *id})
}
