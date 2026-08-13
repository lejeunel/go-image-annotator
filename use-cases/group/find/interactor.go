package find

import (
	"context"
	"fmt"
)

type Interactor struct {
	Repo
}

func (i *Interactor) Execute(ctx context.Context, name string, out OutputPort) {
	errCtx := "fetching group"
	found, err := i.Repo.Find(name)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.SuccessFindGroup(*found)
}

type Option func(*Interactor)

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		Repo: r,
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
