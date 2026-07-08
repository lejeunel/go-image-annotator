package modify_bbox

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
	annotationId, err := a.NewAnnotationIdFromString(r.AnnotationId)
	if err != nil {
		out.Error(err)
		return
	}
	group, err := i.annotationRepo.GroupOfAnnotation(*annotationId)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching annotation group: %w", errCtx, err))
		return
	}

	if group != nil {
		if err := i.auth.Annotate(ctx, *group); err != nil {
			out.Error(fmt.Errorf("%v: authenticating: %w", errCtx, err))
			return
		}
	}
	label, err := i.findLabel(r.Label)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching label %v: %w", errCtx, r.Label, err))
		return
	}
	u, err := i.validate(r.Xc, r.Yc, r.Width, r.Height, *label, r.Angle)
	if err != nil {
		out.Error(fmt.Errorf("%v: validating coordinates: %w", errCtx, err))
		return
	}

	if err := i.update(ctx, *annotationId, *u); err != nil {
		out.Error(fmt.Errorf("%v: updating: %w", errCtx, err))
		return
	}
	out.SuccessUpdateBox(Response{})

}
func (i Interactor) update(ctx context.Context, id a.AnnotationId, upd a.BoundingBoxUpdatables) error {
	var userId *u.UserId
	user := u.IdentityFromContext(ctx)
	if user != nil {
		userId = &user.Id
	}
	now := i.clock.Now()

	if err := i.annotationRepo.UpdateBoundingBox(id, upd, userId, &now); err != nil {
		return err
	}
	return nil

}
func (i Interactor) validate(xc float32, yc float32, width float32,
	height float32, label lbl.Label, angle float32) (*a.BoundingBoxUpdatables, error) {

	if err := a.ValidateBoundingBox(xc, yc, width, height, angle); err != nil {
		return nil, err
	}
	return &a.BoundingBoxUpdatables{LabelId: label.Id, Xc: xc, Yc: yc, Width: width, Height: height,
		Angle: angle}, nil

}
func (i Interactor) findLabel(name string) (*lbl.Label, error) {

	label, err := i.labelRepo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil

}
