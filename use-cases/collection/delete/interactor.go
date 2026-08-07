package delete

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jonboulle/clockwork"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	ev "github.com/lejeunel/go-image-annotator/entities/event"
	im "github.com/lejeunel/go-image-annotator/entities/image"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	event_logger "github.com/lejeunel/go-image-annotator/modules/event-logger"
	ims "github.com/lejeunel/go-image-annotator/modules/image-store"
	job_queue "github.com/lejeunel/go-image-annotator/modules/job-queue"
)

type Transactor interface {
	RunInTx(fn func(Repos) error) error
}

type Interactor struct {
	Repos
	Transactor
	ImageStore ims.Interface
	Auth
	JobQueue    job_queue.Interface
	EventLogger event_logger.Interface
	slog.Logger
	clockwork.Clock
}

func New(
	r Repos,
	tra Transactor,
	ims ims.Interface,
	jq job_queue.Interface,
	el event_logger.Interface,
	logger slog.Logger,
	opts ...Option,
) Interactor {
	i := &Interactor{
		Repos:       r,
		Transactor:  tra,
		ImageStore:  ims,
		EventLogger: el,
		Logger:      logger,
		JobQueue:    jq,
		Clock:       clockwork.NewRealClock(),
		Auth:        auth.NewVoidAuth(),
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i Interactor) Execute(ctx context.Context, name string, out OutputPort) {
	errCtx := "deleting collection"

	collection, err := i.CollectionRepo.Find(name)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching collection: %w", errCtx, err))
		return
	}
	if collection.Group != nil {
		if err := i.Auth.DeleteCollection(ctx, *collection.Group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}
	user := u.IdentityFromContext(ctx)
	if user == nil {
		out.Error(fmt.Errorf("%v: failed fetching user id from context", errCtx))
		return
	}

	task := t.NewTask(t.NewTaskId(), user.Id, t.CollectionDeleteTask)
	if err := i.EventLogger.InitTask(
		task.Id, task.Type, task.Issuer); err != nil {
		out.Error(fmt.Errorf("%v: pushing init task to logger: %w", errCtx, err))
		return
	}
	if err := i.EventLogger.AddEvent(task.Id, ev.Event{Time: i.Clock.Now(), State: ev.PendingTask}); err != nil {
		out.Error(fmt.Errorf("%v: adding pending status: %w", errCtx, err))
		return
	}

	i.JobQueue.Submit(func() {
		i.runTask(task, *collection)
	})
	out.SuccessDeleteCollection(Response{Id: task.Id, Type: task.Type, Issuer: user.Id})
}

func (i *Interactor) runTask(task t.Task, collection clc.Collection) {
	errCtx := fmt.Errorf("running delete collection task")
	i.Logger.Info(fmt.Sprintf("started delete task %v", task.Id))

	extra := map[string]string{"collection": collection.Name}
	if err := i.EventLogger.AddEvent(task.Id, ev.Event{Time: i.Clock.Now(), State: ev.StartedTask, Extra: extra}); err != nil {
		i.Logger.Error(
			fmt.Errorf("%w: logging event upon delete task startup: %w", errCtx, err).Error(),
		)
		return
	}

	for baseImage, err := range i.ImageRepo.Iterate(im.Filtering{Collection: &collection.Name}, 1) {
		if err := i.Transactor.RunInTx(func(tx Repos) error {
			if err != nil {
				return fmt.Errorf("%w: iterating on images: %w", errCtx, err)
			}
			if err := tx.AnnotationRepo.RemoveAllAnnotations(baseImage.ImageId, collection.Name); err != nil {
				return fmt.Errorf("%w: deleting annotations: %w", errCtx, err)
			}
			if err := tx.ImageRepo.RemoveImageFromCollection(baseImage.ImageId, collection.Id); err != nil {
				return fmt.Errorf("%w: adding image to collection: %w", errCtx, err)
			}

			isUsed, err := tx.ImageRepo.IsUsed(baseImage.ImageId)
			if err != nil {
				return fmt.Errorf("%w: checking whether image %v is used in another collection: %w", errCtx, baseImage.ImageId, err)
			}
			if !*isUsed {
				i.ImageStore.DeleteAsset(baseImage.ImageId)
			}
			return nil
		}); err != nil {
			i.EventLogger.AddEvent(
				task.Id,
				ev.Event{Time: i.Clock.Now(), State: ev.FailedTask, Error: err.Error()},
			)
			i.Logger.Error(err.Error())
			return
		}
	}

	if err := i.CollectionRepo.Delete(collection.Name); err != nil {
		i.Logger.Error(fmt.Errorf("%w: deleting collection: %w", errCtx, err).Error())
		return
	}

	i.EventLogger.AddEvent(task.Id, ev.Event{Time: i.Clock.Now(), State: ev.DoneTask})
}

func (i *Interactor) LogError(id t.TaskId, err error) {
	i.EventLogger.AddEvent(
		id,
		ev.Event{Time: i.Clock.Now(), State: ev.FailedTask, Error: err.Error()},
	)
	i.Logger.Error(err.Error())
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}
