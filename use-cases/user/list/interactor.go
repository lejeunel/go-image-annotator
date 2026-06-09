package list

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lejeunel/go-image-annotator/shared/auth"
	"github.com/lejeunel/go-image-annotator/shared/logging"
	"github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	repo   Repo
	logger *slog.Logger
	auth   Auth
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	if err := i.auth.ListUsers(ctx); err != nil {
		i.handleError(err, out)
		return
	}
	found, err := i.repo.List(r)
	if err != nil {
		i.handleError(err, out)
		return
	}

	count, err := i.repo.Count()
	if err != nil {
		i.handleError(err, out)
		return
	}

	response := Response{
		Pagination: pagination.New(int64(r.Page), r.PageSize, count),
	}
	for _, f := range found {
		response.Users = append(response.Users, UserResponse{Id: f.Id, Groups: f.Groups, Roles: f.Roles})
	}
	out.Success(response)

}

func (i *Interactor) handleError(err error, out OutputPort) {
	errCtx := "listing users"
	err = fmt.Errorf("%v: %w", errCtx, err)
	i.logger.Error(errCtx, "error", err)
	out.Error(err)
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.auth = a
	}
}

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r, logger: logging.NewNoOpLogger(),
		auth: auth.PassThroughAuth{}}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
