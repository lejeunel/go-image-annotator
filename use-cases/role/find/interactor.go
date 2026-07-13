package find

import (
	"context"
	"fmt"
)

type Interactor struct {
	repo Repo
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	errCtx := "fetching role"
	found, err := i.repo.Find(r.Name)
	if err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	out.Success(Response{Name: found.Name, Description: found.Description})

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
