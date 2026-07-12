package find

import (
	"context"
	"fmt"
)

type Interactor struct {
	repo Repo
}

func (i *Interactor) Execute(ctx context.Context, name string, out OutputPort) {
	found, err := i.repo.FindLabel(name)
	if err != nil {
		out.Error(fmt.Errorf("fetching label by name %v: %w", name, err))
		return
	}

	out.SuccessFindLabel(*found)

}

type Option func(*Interactor)

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
