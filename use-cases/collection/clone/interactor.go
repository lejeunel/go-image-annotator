package clone

import (
	"context"
	"fmt"

	t "github.com/lejeunel/go-image-annotator/entities/task"
	u "github.com/lejeunel/go-image-annotator/entities/user"
)

type Cloner interface {
	Clone(t.CloneTask) error
}

type Interactor struct {
	Cloner
	Auth
}

func New(c Cloner, opts ...Option) Interactor {
	i := &Interactor{Cloner: c}
	for _, opt := range opts {
		opt(i)
	}
	return *i
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

	if r.DestinationGroup != nil {
		if err := i.Auth.CloneCollection(ctx, *r.DestinationGroup); err != nil {
			out.Error(fmt.Errorf("%v: %w", errCtx, err))
			return
		}
	}

	specs := t.NewSpecs(t.NewTaskId(), user.Id, t.CollectionCloneTask)
	task := t.NewCloneTask(t.NewTaskId(), user.Id,
		r.Source, r.Destination, t.WithDeepClone())
	if err := i.Cloner.Clone(task); err != nil {
		out.Error(fmt.Errorf("%w: %w", errCtx, err))
		return
	}

	out.SuccessSubmitCloneTask(specs)
}
