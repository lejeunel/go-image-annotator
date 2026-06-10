package read

import (
	"context"
	"fmt"
)

type Interactor struct {
	repo Repo
}

func (i *Interactor) Execute(ctx context.Context, r Request, out OutputPort) {
	found, err := i.repo.FindLabel(r.Name)
	if err != nil {
		out.Error(fmt.Errorf("fetching label by name %v: %w", r.Name, err))
		return
	}

	out.Success(Response{Name: found.Name, Description: found.Description})

}

type Option func(*Interactor)

func New(r Repo, opts ...Option) *Interactor {
	i := &Interactor{repo: r}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
