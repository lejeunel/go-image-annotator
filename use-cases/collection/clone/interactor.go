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
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	jq "github.com/lejeunel/go-image-annotator/modules/job-queue"
)

type ImageStore interface {
	Find(base im.BaseImage) (*im.Image, error)
	Copy(src clc.CollectionName, id im.ImageId, dst clc.CollectionName, deep bool) error
}

type Interactor struct {
	ImageStore
	ImageRepo
	CollectionRepo
	GroupRepo
	el.IEventLogger
	Auth
	clockwork.Clock
	slog.Logger
	jq.JobQueue
}

func New(ims ImageStore, ir ImageRepo, c CollectionRepo, g GroupRepo,
	l el.IEventLogger, logger slog.Logger, j jq.JobQueue,
	opts ...Option,
) Interactor {
	itr := &Interactor{ims, ir, c, g, l, auth.NewVoidAuth(), clockwork.NewRealClock(), logger, j}
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

	if err := i.IEventLogger.InitTask(
		task.Id, task.Type, task.Issuer); err != nil {
		out.Error(fmt.Errorf("%v: pushing init task to logger: %w", errCtx, err))
		return
	}
	if err := i.IEventLogger.AddEvent(
		task.Id,
		e.Event{Time: i.Clock.Now(), State: e.PendingTask},
	); err != nil {
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
		errs = errors.Join(
			errs,
			fmt.Errorf("destination collection %q already exists", destination),
		)
	}
	return errs
}

func (i *Interactor) LogError(id t.TaskId, err error) {
	i.IEventLogger.AddEvent(
		id,
		e.Event{Time: i.Clock.Now(), State: e.FailedTask, Error: err.Error()},
	)
	i.Logger.Error(err.Error())
}

func (i *Interactor) runTask(
	task t.Task,
	source string,
	destination string,
	group *grp.Group,
	deep bool,
) {
	errCtx := fmt.Errorf("running collection cloning task")
	i.Logger.Info(fmt.Sprintf("started clone task %v", task.Id))

	if err := i.checkCollections(source, destination); err != nil {
		i.LogError(task.Id, fmt.Errorf("%w: %w", errCtx, err))
		return
	}
	extra := map[string]string{
		"source-collection":      source,
		"destination-collection": destination,
		"deep-copy":              strconv.FormatBool(deep),
	}
	if err := i.IEventLogger.AddEvent(
		task.Id,
		e.Event{Time: i.Clock.Now(), State: e.StartedTask, Extra: extra},
	); err != nil {
		i.Logger.Error(
			fmt.Errorf("%w: logging event upon cloning task startup: %w", errCtx, err).Error(),
		)
		return
	}

	dst := clc.NewCollection(clc.NewCollectionId(), destination, clc.WithCreatedAt(i.Clock.Now()))
	if group != nil {
		dst.Group = &group.Name
	}

	if err := i.CollectionRepo.Create(dst); err != nil {
		i.IEventLogger.AddEvent(
			task.Id,
			e.Event{Time: i.Clock.Now(), State: e.FailedTask, Error: err.Error()},
		)
		i.Logger.Error(err.Error())
		return
	}

	for baseImage, err := range i.ImageRepo.Iterate(im.Filtering{Collection: &source}, 1) {
		if err != nil {
			i.LogError(task.Id, err)
			return
		}
		if err := i.ImageStore.Copy(source, baseImage.ImageId, dst.Name, deep); err != nil {
			i.LogError(task.Id, err)
			return

		}
	}
	i.IEventLogger.AddEvent(task.Id, e.Event{Time: i.Clock.Now(), State: e.DoneTask})
}
