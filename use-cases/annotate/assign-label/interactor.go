package assign_label

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"

	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	sauth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
)

type Interface interface {
	Execute(ctx context.Context, r Request, out OutputPort)
}
type ImageStore interface {
	Find(base im.BaseImage) (*im.Image, error)
}

type Interactor struct {
	AnnotationRepo
	LabelRepo
	ImageStore
	auth.Auth
	clockwork.Clock
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "assigning label to image"

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
		if err := i.Auth.Annotate(ctx, *image.Collection.Group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	label, err := i.findLabel(r.Label)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	imageLabel, err := i.addLabel(ctx, image.Id, image.Collection.Name, *label)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessAddLabel(Response{
		ImageId:      r.ImageId,
		Collection:   r.Collection,
		Label:        r.Label,
		AnnotationId: imageLabel.Id.String(),
	})
}

func (i Interactor) findLabel(name string) (*lbl.Label, error) {
	label, err := i.LabelRepo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil
}

func (i Interactor) findImage(imageId im.ImageId, collection string) (*im.Image, error) {
	image, err := i.ImageStore.Find(im.BaseImage{ImageId: imageId, Collection: collection})
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (i Interactor) addLabel(
	ctx context.Context,
	imageId im.ImageId,
	collection clc.CollectionName,
	label lbl.Label,
) (*an.ImageLabel, error) {
	var userId *u.UserId
	user := u.IdentityFromContext(ctx)
	if user != nil {
		userId = &user.Id
	}
	now := i.Clock.Now()

	imageLabel := an.NewImageLabel(label)
	if err := i.AnnotationRepo.AddImageLabel(
		imageId,
		collection,
		imageLabel,
		userId,
		&now,
	); err != nil {
		return nil, err
	}
	return &imageLabel, nil
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

func New(repo AnnotationRepo, labelRepo LabelRepo, store ImageStore, opts ...Option) Interactor {
	i := &Interactor{
		AnnotationRepo: repo,
		LabelRepo:      labelRepo,
		ImageStore:     store,
		Clock:          clockwork.NewRealClock(),

		Auth: sauth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
