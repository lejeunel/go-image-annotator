package add_polygon

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"
	a "github.com/lejeunel/go-image-annotator/entities/annotation"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	sauth "github.com/lejeunel/go-image-annotator/modules/auth"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
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
		auth:           sauth.NewVoidAuth()}
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
	errCtx := "adding polygon"

	imageId, err := im.NewImageIdFromString(r.ImageId)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	image, err := i.findImage(imageId, r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if image.Collection.Group != nil {
		if err := i.auth.Annotate(ctx, image.Collection.Group.Name); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	label, err := i.findLabel(r.Label)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	poly := a.NewPolygon(a.NewAnnotationId(), r.Points, *label)
	if err := i.addPolygon(ctx, image, poly); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessAddPolygon(Response{poly.Id})

}
func (i Interactor) addPolygon(ctx context.Context, image *im.Image, poly a.Polygon) error {
	var userId *u.UserId
	user := u.IdentityFromContext(ctx)
	if user != nil {
		userId = &user.Id
	}
	now := i.clock.Now()
	if err := i.annotationRepo.AddPolygon(image.Id, image.Collection.Id, poly, userId, &now); err != nil {
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
func (i Interactor) findImage(imageId im.ImageId, collectionName string) (*im.Image, error) {
	image, err := i.imageStore.Find(im.BaseImage{ImageId: imageId, Collection: collectionName})
	if err != nil {
		return nil, err
	}
	return image, nil
}
