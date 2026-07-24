package find

import (
	"context"
	"fmt"
)

type Interactor struct {
	repo Repo
}

func (i *Interactor) Execute(ctx context.Context, name string, out OutputPort) {
	errCtx := "fetching role"
	found, err := i.repo.Find(name)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.SuccessFindRole(*found)

}

type Option func(*Interactor)

func New(r Repo, opts ...Option) Interactor {
	i := &Interactor{
		repo: r,
	}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
