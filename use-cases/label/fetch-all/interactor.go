package fetchall

import (
	"context"
	"fmt"
	"log/slog"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

var defaultLabelCountLimit = 200

type Interface interface {
	Execute(context.Context, OutputPort)
}
type Interactor struct {
	repo       Repo
	logger     *slog.Logger
	countLimit int
}

func (i Interactor) Execute(ctx context.Context, out OutputPort) {
	errCtx := "listing label"
	count, err := i.repo.Count()
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	if count > int64(i.countLimit) {
		out.Error(fmt.Errorf("checking whether current label count (%v) exceeds limit (%v): %w",
			count, i.countLimit, e.ErrLabelLimitExceeded))
		return
	}

	labels, err := i.repo.FetchAll()
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.SuccessFetchLabels(Response{labels})

}

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r,
		countLimit: defaultLabelCountLimit}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

type Option func(*Interactor)

func WithLimit(limit int) Option {
	return func(c *Interactor) {
		c.countLimit = limit
	}
}
