package add_bbox

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
	sauth "github.com/lejeunel/go-image-annotator/shared/auth"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
)

type Interface interface {
	Execute(context.Context, Request, OutputPort)
}

type Interactor struct {
	imageStore     st.Interface
	annotationRepo Repo
	labelRepo      LabelRepo
	auth           auth.Auth
	clock          clockwork.Clock
}

func New(imageStore st.Interface, repo Repo, labelRepo LabelRepo, opts ...Option) Interactor {
	i := &Interactor{
		annotationRepo: repo,
		labelRepo:      labelRepo,
		imageStore:     imageStore,
		clock:          clockwork.NewRealClock(),
		auth:           sauth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return *i
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

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "adding bounding box"
	image, err := i.findImage(r.ImageId, r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if image.Collection.Group != nil {
		if err := i.auth.AnnotateGroup(ctx, image.Collection.Group.Name); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	label, err := i.findLabel(r.Label)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	box := a.NewBoundingBox(a.NewAnnotationId(), r.Xc, r.Yc, r.Width, r.Height, *label,
		a.WithAngle(r.Angle))
	if err := i.validateBox(image, box); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.addBox(ctx, image, box); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessAddBox(Response{box.Id})

}
func (i Interactor) addBox(ctx context.Context, image *im.Image, box a.BoundingBox) error {
	var userId *u.UserId
	user := ip.IdentityFromContext(ctx)
	if user != nil {
		userId = &user.Id
	}
	now := i.clock.Now()
	if err := i.annotationRepo.AddBoundingBox(image.Id, image.Collection.Id, box, userId, &now); err != nil {
		return err
	}
	return nil
}
func (i Interactor) validateBox(image *im.Image, box a.BoundingBox) error {
	if err := image.AddBoundingBox(box); err != nil {
		return err
	}
	return nil
}
func (i Interactor) findLabel(name string) (*lbl.Label, error) {
	label, err := i.labelRepo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil
}
func (i Interactor) findImage(imageId string, collectionName string) (*im.Image, error) {
	image, err := i.imageStore.Find(im.BaseImage{ImageId: imageId, Collection: collectionName})
	if err != nil {
		return nil, err
	}
	return image, nil
}
