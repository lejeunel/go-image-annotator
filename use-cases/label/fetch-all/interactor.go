package fetchall

import (
	"context"
	"fmt"
	"log/slog"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/shared/logging"
)

var defaultLabelCountLimit = 200

type Interface interface {
	Execute(context.Context, OutputPort)
}
type Interactor struct {
	repo       Repo
	logger     *slog.Logger
	countLimit int
	auth       Auth
}

func (i *Interactor) Execute(ctx context.Context, out OutputPort) {
	if err := i.auth.FetchAllLabels(ctx); err != nil {
		i.handleError(err, out)
		return
	}
	count, err := i.repo.Count()
	if err != nil {
		i.handleError(e.ErrInternal, out)
		return
	}
	if count > int64(i.countLimit) {
		i.handleError(fmt.Errorf("checking whether current label count (%v) exceeds limit (%v): %w",
			count, i.countLimit, e.ErrLabelLimitExceeded), out)
		return
	}

	labels, err := i.repo.FetchAll()
	if err != nil {
		i.handleError(err, out)
		return
	}
	out.SuccessFetchLabels(Response{labels})

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "listing label"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

func NewInteractor(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r, logger: logging.NewNoOpLogger(),
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

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}
