package find

import (
	"context"
	"fmt"
)

type Interactor struct {
	repo Repo
}

func (i Interactor) Execute(ctx context.Context, name string, out OutputPort) {
	errCtx := "fetching collection"
	found, err := i.repo.FindCollectionByName(name)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}
	out.SuccessFindCollection(*found)
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
