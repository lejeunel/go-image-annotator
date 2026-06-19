package update_label

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	sauth "github.com/lejeunel/go-image-annotator/modules/auth"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	annotationRepo AnnotationRepo
	labelRepo      LabelRepo
	auth           auth.Auth
	clock          clockwork.Clock
}

type Option func(*Interactor)

func WithAuth(a auth.Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func WithClock(c clockwork.Clock) Option {
	return func(i *Interactor) {
		i.clock = c
	}
}

func New(repo AnnotationRepo, labelRepo LabelRepo, opts ...Option) Interactor {
	i := &Interactor{annotationRepo: repo,
		labelRepo: labelRepo,
		clock:     clockwork.NewRealClock(),
		auth:      sauth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "updating bounding box properties"
	label, err := i.labelRepo.FindLabel(r.Label)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	id, err := a.NewAnnotationIdFromString(r.AnnotationId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	group, err := i.annotationRepo.GroupOfAnnotation(*id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if group != nil {
		if err := i.auth.Annotate(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	var userId *u.UserId
	user := u.IdentityFromContext(ctx)
	if user != nil {
		userId = &user.Id
	}
	now := i.clock.Now()

	err = i.annotationRepo.UpdateLabelOfAnnotation(*id, label.Id, userId, &now)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessUpdateLabel(Response{})

}
