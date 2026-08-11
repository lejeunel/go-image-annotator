package ingest

import (
	"context"
	"fmt"
	"github.com/jonboulle/clockwork"
	"io"
	"log/slog"

	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	ev "github.com/lejeunel/go-image-annotator/entities/event"
	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	aig "github.com/lejeunel/go-image-annotator/modules/archive-ingester"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	el "github.com/lejeunel/go-image-annotator/modules/event-logger"
	iig "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	jq "github.com/lejeunel/go-image-annotator/modules/job-queue"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Auth interface {
	IngestImage(ctx context.Context, group string) error
}

type ImageIngester interface {
	Ingest(iig.Request) (*iig.Response, error)
}

type ArchiveIngester interface {
	IngestArchive(aig.Request) (aig.Response, error)
}

type TemporaryFileStore interface {
	Store(string, io.Reader) error
	GetReaderAt(string) (io.ReaderAt, int64, error)
	Delete(string) error
}

type Interactor struct {
	ArchiveIngester
	Auth
	CollectionRepo
	TemporaryFileStore
	el.IEventLogger
	clockwork.Clock
	jq.JobQueue
	slog.Logger
}
type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(aig ArchiveIngester,
	cr CollectionRepo,
	tfs TemporaryFileStore,
	el el.IEventLogger,
	logger slog.Logger,
	jq jq.JobQueue,
	opts ...Option) Interactor {

	i := &Interactor{
		ArchiveIngester:    aig,
		CollectionRepo:     cr,
		TemporaryFileStore: tfs,
		Auth:               auth.NewVoidAuth(),
		IEventLogger:       el,
		Clock:              clockwork.NewRealClock(),
		JobQueue:           jq,
		Logger:             logger,
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}

func (i Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := fmt.Errorf("ingesting image archive")
	collection, err := i.findCollectionByName(r.Collection)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	if collection.Group != nil {
		if err := i.Auth.IngestImage(ctx, *collection.Group); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	user := u.IdentityFromContext(ctx)
	if user == nil {
		out.Error(
			fmt.Errorf(
				"%w: extracting user identity failed from context: %w",
				errCtx,
				e.ErrAuthentication,
			),
		)
		return
	}
	task := t.NewTask(t.NewTaskId(), user.Id, t.IngestArchiveTask)
	tmpFileName := fmt.Sprintf("%v.zip", task.Id)
	if err := i.TemporaryFileStore.Store(tmpFileName, r.Reader); err != nil {
		out.Error(
			fmt.Errorf("%w: storing archive in temporary location: %w", errCtx, err),
		)
		return
	}
	if err := i.IEventLogger.InitTask(task.Id, task.Type, task.Issuer); err != nil {
		out.Error(
			fmt.Errorf("%w: initializing ingestion task: %w", errCtx, err),
		)
		return
	}
	if err := i.IEventLogger.AddEvent(
		task.Id,
		ev.Event{Time: i.Clock.Now(), State: ev.PendingTask},
	); err != nil {
		out.Error(fmt.Errorf("%v: adding pending status: %w", errCtx, err))
		return
	}

	i.JobQueue.Submit(func() {
		i.runTask(task, user.Id, r.Collection, tmpFileName)
	})
	out.SuccessSubmitIngestArchiveTask(Response{Id: task.Id, Issuer: task.Issuer, Type: task.Type})
}

func (i Interactor) runTask(task t.Task, user u.UserId, collection clc.CollectionName, filename string) {
	reader, size, err := i.TemporaryFileStore.GetReaderAt(filename)
	if err != nil {
		i.LogError(task.Id, fmt.Errorf("ingesting archive: reading archive from temporary store: %w", err))
		return
	}
	i.IEventLogger.AddEvent(task.Id,
		ev.Event{
			Time:  i.Clock.Now(),
			State: ev.StartedTask,
			Extra: map[string]string{"collection": collection},
		})
	resp, err := i.ArchiveIngester.IngestArchive(aig.Request{UserId: user, Collection: collection,
		ReaderAt: reader, Size: size,
	})
	if err != nil {
		i.LogError(task.Id, err)
		if err := i.TemporaryFileStore.Delete(filename); err != nil {
			i.LogError(task.Id, fmt.Errorf("ingesting archive: deleting file from temporary store: %w", err))
			return
		}
		return
	}

	i.IEventLogger.AddEvent(
		task.Id,
		ev.Event{Time: i.Clock.Now(), State: ev.DoneTask,
			Extra: map[string]string{"num-ingested-images": fmt.Sprintf("%v", len(resp.ImageIds))}},
	)

}

func (i *Interactor) LogError(id t.TaskId, err error) {
	i.IEventLogger.AddEvent(
		id,
		ev.Event{Time: i.Clock.Now(), State: ev.FailedTask, Error: err.Error()},
	)
	i.Logger.Error(err.Error())
}

func (i Interactor) findCollectionByName(name string) (*clc.Collection, error) {
	collection, err := i.CollectionRepo.Find(name)
	baseErr := fmt.Errorf("finding collection with name %v", name)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", baseErr, err)
	}
	return collection, nil
}
