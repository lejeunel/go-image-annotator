package clone

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jonboulle/clockwork"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	e "github.com/lejeunel/go-image-annotator/entities/event"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	event_logger "github.com/lejeunel/go-image-annotator/modules/event-logger"
	st "github.com/lejeunel/go-image-annotator/modules/image-store"
)

type JobQueue interface {
	Submit(func())
}

type Interactor struct {
	ImageRepo
	CollectionRepo
	AnnotationRepo
	GroupRepo
	Store       st.Interface
	EventLogger event_logger.Interface
	Auth
	clockwork.Clock
	slog.Logger
	JobQueue
}

func New(i ImageRepo, c CollectionRepo, a AnnotationRepo, g GroupRepo,
	s st.Interface, l event_logger.Interface, logger slog.Logger, j JobQueue,
	opts ...Option) Interactor {
	itr := &Interactor{i, c, a, g, s, l, auth.NewVoidAuth(), clockwork.NewRealClock(), logger, j}
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

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("initiating cloning collection task")
	user := u.IdentityFromContext(ctx)
	if user == nil {
		out.Error(fmt.Errorf("%w: failed fetching user id from context", errCtx))
		return
	}

	task := t.NewTask(t.NewTaskId(), user.Id, t.CollectionCloneTask)

	var group *grp.Group
	if r.DestinationGroup != nil {
		if err := i.Auth.CloneCollection(ctx, *r.DestinationGroup); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
		var err error
		group, err = i.GroupRepo.Find(*r.DestinationGroup)
		if err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	if err := i.checkCollections(r.Source, r.Destination); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if err := i.EventLogger.InitTask(
		task.Id, task.Type, task.Issuer,
		"source-collection", r.Source,
		"destination-collection", r.Destination,
		"deep-copy", strconv.FormatBool(r.Deep)); err != nil {
		out.Error(fmt.Errorf("%v: pushing init task to logger: %w", errCtx, err))
		return
	}

	if err := i.EventLogger.AddEvent(i.Clock.Now(), task.Id, e.PendingTask, nil); err != nil {
		out.Error(fmt.Errorf("%v: adding pending status: %w", errCtx, err))
		return
	}

	i.JobQueue.Submit(func() {
		i.runTask(task, r.Source, r.Destination, group, r.Deep)
	})
	out.SuccessSubmitCloneTask(Response{Id: task.Id, Issuer: task.Issuer, Type: task.Type})
}

func (i *Interactor) checkCollections(source, destination string) error {
	existsSrc, errSrc := i.CollectionRepo.Exists(source)
	existsDst, errDst := i.CollectionRepo.Exists(destination)

	if err := errors.Join(errSrc, errDst); err != nil {
		return fmt.Errorf("checking source and destination collections existence: %w", err)
	}

	var errs error
	if !existsSrc {
		errs = errors.Join(errs, fmt.Errorf("source collection %q does not exist", source))
	}
	if existsDst {
		errs = errors.Join(errs, fmt.Errorf("destination collection %q already exists", destination))
	}
	return errs

}

func (i *Interactor) LogError(id t.TaskId, err error) {
	i.EventLogger.AddEvent(i.Clock.Now(), id, e.FailedTask, err)
	i.Logger.Error(err.Error())
}

func (i *Interactor) runTask(task t.Task, source string, destination string, group *grp.Group, deep bool) {
	errCtx := fmt.Errorf("running collection cloning task")

	if err := i.checkCollections(source, destination); err != nil {
		i.LogError(task.Id, fmt.Errorf("%w: %w", errCtx, err))
		return
	}

	if err := i.EventLogger.AddEvent(i.Clock.Now(), task.Id, e.StartedTask, nil); err != nil {
		i.Logger.Error(fmt.Errorf("%w: logging event upon cloning task startup: %w", errCtx, err).Error())
		return
	}

	var dst clc.Collection
	if group != nil {
		dst = clc.NewCollection(clc.NewCollectionId(), destination, clc.WithGroup(*group))
	} else {
		dst = clc.NewCollection(clc.NewCollectionId(), destination)
	}
	if err := i.CollectionRepo.Create(dst); err != nil {
		i.EventLogger.AddEvent(i.Clock.Now(), task.Id, e.FailedTask, err)
		i.Logger.Error(err.Error())
		return
	}

	for baseImage, err := range i.ImageRepo.Iterate(im.Filtering{Collection: &source}, 1) {
		if err != nil {
			err = fmt.Errorf("%w: iterating on images: %w", errCtx, err)
			i.EventLogger.AddEvent(i.Clock.Now(), task.Id, e.FailedTask, err)
			i.Logger.Error(err.Error())
			return
		}
		image, err := i.Store.Find(baseImage)
		if err != nil {
			i.LogError(task.Id, fmt.Errorf("%w: finding source image: %w", errCtx, err))
			return
		}
		if err := i.ImageRepo.AddToCollection(image.Id, dst.Id); err != nil {
			i.LogError(task.Id, fmt.Errorf("%w: adding image to collection: %w", errCtx, err))
			return
		}

		if deep {
			for _, label := range image.Labels {
				if err := i.AnnotationRepo.AddImageLabel(image.Id, image.Collection.Id, label, label.Author, label.Time); err != nil {
					i.LogError(task.Id, fmt.Errorf("%w: adding image label: %w", errCtx, err))
					return
				}
			}

			for _, box := range image.BoundingBoxes {
				if err := i.AnnotationRepo.AddBoundingBox(image.Id, image.Collection.Id, box, box.Author, box.Time); err != nil {
					i.LogError(task.Id, fmt.Errorf("%w: adding bounding boxes: %w", errCtx, err))
					return
				}
			}
			for _, poly := range image.Polygons {
				if err := i.AnnotationRepo.AddPolygon(image.Id, image.Collection.Id, poly, poly.Author, poly.Time); err != nil {
					i.LogError(task.Id, fmt.Errorf("%w: adding polygons: %w", errCtx, err))
					return
				}
			}

		}
	}
	i.EventLogger.AddEvent(i.Clock.Now(), task.Id, e.DoneTask, nil)
}
