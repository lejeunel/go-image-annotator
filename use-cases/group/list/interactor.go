package list

import (
	"context"
	"fmt"
)

type Interactor struct {
	repo Repo
}

func (i *Interactor) Execute(ctx context.Context, out OutputPort) {
	errCtx := "listing groups"

	found, err := i.repo.List()
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessListGroups(found)
}

type Option func(*Interactor)

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{repo: r}

	for _, opt := range opts {
		opt(i)
	}
	return *i
}
