package list

import (
	"context"
	"fmt"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type Interactor struct {
	logger TaskLogger
}

func (i Interactor) Execute(ctx context.Context, r pa.PaginationParams, out OutputPort) {
	errCtx := "listing tasks"
	if err := pa.Validate(r.Page, r.PageSize); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	user := u.IdentityFromContext(ctx)
	if user == nil {
		out.Error(fmt.Errorf("%v: fetching user identity: %w", errCtx, e.ErrAuthentication))
		return
	}

	found, err := i.logger.ListUserTasks(user.Id, r)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	count, err := i.logger.Count(user.Id)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	response := Response{Pagination: pa.New(int64(r.Page), r.PageSize, *count)}
	for _, f := range found {
		response.Tasks = append(response.Tasks, f)
	}
	out.SuccessListTasks(response)
}

type Option func(*Interactor)

func New(r TaskLogger, opts ...Option) Interactor {
	i := &Interactor{logger: r}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
