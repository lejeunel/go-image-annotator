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
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
	sauth "github.com/lejeunel/go-image-annotator/shared/auth"
	ip "github.com/lejeunel/go-image-annotator/shared/identity_provider"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"github.com/lejeunel/go-image-annotator/use-cases/annotate/auth"
	"log/slog"
)

type Interface interface {
	Execute(ctx context.Context, r Request, out OutputPort)
}

type Interactor struct {
	repo   Repo
	store  st.Interface
	logger *slog.Logger
	auth   auth.Auth
	clock  clockwork.Clock
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {

	image, err := i.findImage(r.ImageId, r.Collection)
	if err != nil {
		i.handleError(err, out)
		return
	}

	if image.Collection.Group != nil {
		if err := i.auth.AnnotateGroup(ctx, image.Collection.Group.Name); err != nil {
			i.handleError(err, out)
			return
		}

	}

	label, err := i.findLabel(r.Label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	imageLabel, err := i.addLabel(ctx, image.Id, image.Collection.Id, *label)
	if err != nil {
		i.handleError(err, out)
		return
	}

	out.SuccessAddLabel(Response{
		ImageId:      r.ImageId,
		Collection:   r.Collection,
		Label:        r.Label,
		AnnotationId: imageLabel.Id.String()})
}
func (i Interactor) handleError(err error, out OutputPort) {
	errCtx := "assigning label to image"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)

}
func (i Interactor) findLabel(name string) (*lbl.Label, error) {
	label, err := i.repo.FindLabel(name)
	if err != nil {
		return nil, err
	}
	return label, nil

}
func (i Interactor) findImage(imageId string, collection string) (*im.Image, error) {
	image, err := i.store.Find(im.BaseImage{ImageId: imageId, Collection: collection})
	if err != nil {
		return nil, err
	}
	return image, nil
}
func (i Interactor) addLabel(ctx context.Context, imageId im.ImageId, collectionId clc.CollectionId, label lbl.Label) (*an.ImageLabel, error) {
	var userId *u.UserId
	user := ip.IdentityFromContext(ctx)
	if user != nil {
		userId = &user.Id
	}
	now := i.clock.Now()

	imageLabel := an.NewImageLabel(label)
	if err := i.repo.AddImageLabel(imageId, collectionId, imageLabel, userId, &now); err != nil {
		return nil, err
	}
	return &imageLabel, nil

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

func New(repo Repo, store st.Interface, opts ...Option) Interactor {
	i := &Interactor{repo: repo, store: store, logger: logging.NewNoOpLogger(),
		clock: clockwork.NewRealClock(),

		auth: sauth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
