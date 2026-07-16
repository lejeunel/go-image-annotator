package clone

import (
	"context"
	"fmt"
	"strconv"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/entities/event"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	event_logger "github.com/lejeunel/go-image-annotator/modules/event-logger"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
	er "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Interactor struct {
	ImageRepo
	CollectionRepo
	AnnotationRepo
	GroupRepo
	Store       st.Interface
	EventLogger event_logger.Interface
	Auth
}

func New(i ImageRepo, c CollectionRepo, a AnnotationRepo, g GroupRepo, s st.Interface, l event_logger.Interface, opts ...Option) Interactor {
	itr := &Interactor{i, c, a, g, s, l, auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(itr)
	}
	return *itr
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("initiating cloning collection task")
	user := u.IdentityFromContext(ctx)
	if user == nil {
		out.Error(fmt.Errorf("%w: failed fetching user id from context", errCtx))
		return
	}

	task := t.NewCloneTask(t.NewTaskId(), user.Id,
		r.Source, r.Destination, t.WithDeepClone())

	if r.DestinationGroup != nil {
		if err := i.Auth.CloneCollection(ctx, *r.DestinationGroup); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
		task.Group = r.DestinationGroup
	}

	out.SuccessSubmitCloneTask(Response{Id: task.TaskId, Issuer: task.Issuer, Type: task.Type(), State: task.State})
}

func (i Interactor) makeTaskFunc(task t.CloneTask) func() {
	extra := make(map[string]string)
	deep := task.Deep
	extra["source"] = task.Source
	extra["destination"] = task.Destination
	extra["deep"] = strconv.FormatBool(deep)
	event := e.New(task.TaskId, task.Issuer, task.Type())
	event.Extra = extra
	i.EventLogger.Log(event)
	return func() {
		i.EventLogger.Log(event.SetStart())
		errCtx := fmt.Errorf("running collection cloning task")

		exists, err := i.CollectionRepo.Exists(task.Destination)
		if err != nil {
			i.EventLogger.Log(event.SetError(fmt.Errorf("%w: checking whether collection %v exists: %w", errCtx, task.Destination, err)))
			return
		}
		if exists {
			i.EventLogger.Log(event.SetError(fmt.Errorf("%w: checking for duplicate destination collection: %w", errCtx, er.ErrDuplicate)))
			return
		}
		var dst clc.Collection
		group := task.Group
		if group != nil {
			grp, err := i.GroupRepo.Find(*group)
			if err != nil {
				i.EventLogger.Log(event.SetError(fmt.Errorf("%w: checking for group %v: %w", errCtx, *group, er.ErrDuplicate)))
				return
			}
			dst = clc.NewCollection(clc.NewCollectionId(), task.Destination, clc.WithGroup(*grp))
		} else {
			dst = clc.NewCollection(clc.NewCollectionId(), task.Destination)
		}
		if err := i.CollectionRepo.Create(dst); err != nil {
			i.EventLogger.Log(event.SetError(fmt.Errorf("%w: creating collection %v: %w", errCtx, dst.Name, err)))
			return
		}

		for baseImage, err := range i.ImageRepo.Iterate(im.FilteringParams{Collection: &task.Source}, 1) {
			if err != nil {
				i.EventLogger.Log(event.SetError(fmt.Errorf("%w: iterating on images: %w: %w", errCtx, err, er.ErrInternal)))
				return
			}
			if err := i.ImageRepo.AddToCollection(baseImage.ImageId, dst.Id); err != nil {
				i.EventLogger.Log(event.SetError(fmt.Errorf("%w: adding image to collection: %w: %w", errCtx, err, er.ErrInternal)))
				return
			}

			image, err := i.Store.Find(baseImage)
			if err != nil {
				i.EventLogger.Log(event.SetError(fmt.Errorf("%w: %w", errCtx, err)))
				return
			}
			if task.Deep {
				for _, label := range image.Labels {
					if err := i.AnnotationRepo.AddImageLabel(image.Id, image.Collection.Id, label, label.Author, label.Time); err != nil {
						i.EventLogger.Log(event.SetError(fmt.Errorf("%w: adding image label to collection: %w", errCtx, err)))
						return
					}
				}

				for _, box := range image.BoundingBoxes {
					if err := i.AnnotationRepo.AddBoundingBox(image.Id, image.Collection.Id, box, box.Author, box.Time); err != nil {
						i.EventLogger.Log(event.SetError(fmt.Errorf("%w: adding bounding box to collection: %w", errCtx, err)))
						return
					}
				}
				for _, poly := range image.Polygons {
					if err := i.AnnotationRepo.AddPolygon(image.Id, image.Collection.Id, poly, poly.Author, poly.Time); err != nil {
						i.EventLogger.Log(event.SetError(fmt.Errorf("%w: adding polygon to collection: %w", errCtx, err)))
						return
					}
				}

			}
		}
		i.EventLogger.Log(event.SetDone())

	}
}
