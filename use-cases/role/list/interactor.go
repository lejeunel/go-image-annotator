package list

import (
	"context"
	"fmt"
)

type Interactor struct {
	Repo
}

func (i *Interactor) Execute(ctx context.Context, out OutputPort) {
	errCtx := "listing roles"
	found, err := i.Repo.List()
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessListRoles(found)
}

type Option func(*Interactor)

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{Repo: r}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
