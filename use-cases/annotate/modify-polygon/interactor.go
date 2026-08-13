package modify_polygon

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"

	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	sauth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	AnnotationRepo
	LabelRepo
	auth.Auth
	clockwork.Clock
}

type Option func(*Interactor)

func WithAuth(a auth.Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func WithClock(c clockwork.Clock) Option {
	return func(i *Interactor) {
		i.Clock = c
	}
}

func New(repo AnnotationRepo, labelRepo LabelRepo, opts ...Option) Interactor {
	i := &Interactor{
		AnnotationRepo: repo,
		LabelRepo:      labelRepo,
		Clock:          clockwork.NewRealClock(),
		Auth:           sauth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "updating polygon"
	annotationId, err := a.NewAnnotationIdFromString(r.AnnotationId)
	if err != nil {
		out.Error(err)
		return
	}
	group, err := i.AnnotationRepo.GroupOfAnnotation(*annotationId)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching annotation group: %w", errCtx, err))
		return
	}

	if group != nil {
		if err := i.Auth.Annotate(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: authenticating: %w", errCtx, err))
			return
		}
	}
	label, err := i.findLabel(r.Label)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching label %v: %w", errCtx, r.Label, err))
		return
	}
	if err := i.update(
		ctx,
		*annotationId,
		a.PolygonUpdatables{LabelId: label.Id, Points: r.Points},
	); err != nil {
		out.Error(fmt.Errorf("%v: updating: %w", errCtx, err))
		return
	}
	out.SuccessUpdatePolygon(Response{})
}

func (i Interactor) update(ctx context.Context, id a.AnnotationId, upd a.PolygonUpdatables) error {
	var userId *u.UserId
	user := u.IdentityFromContext(ctx)
	if user != nil {
		userId = &user.Id
	}
	now := i.Clock.Now()

	if err := i.AnnotationRepo.UpdatePolygon(id, upd, userId, &now); err != nil {
		return err
	}
	return nil
}

func (i Interactor) findLabel(name string) (*lbl.Label, error) {
	label, err := i.LabelRepo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil
}
