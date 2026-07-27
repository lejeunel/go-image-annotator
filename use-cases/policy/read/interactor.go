package read

import (
	"context"
	"fmt"
	p "github.com/lejeunel/go-image-annotator/entities/policy"
	auth "github.com/lejeunel/go-image-annotator/modules/authorizer"
	"io"
)

type Store interface {
	Get(string) (io.Reader, error)
}

type Auth interface {
	ReadPolicies(ctx context.Context) error
}

type Interactor struct {
	Store
	Auth
}

func (i *Interactor) Execute(ctx context.Context, out OutputPort) {
	errCtx := "fetching policies"
	if err := i.Auth.ReadPolicies(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	r, err := i.Store.Get(p.DefaultPolicyFileName)
	if err != nil {
		out.Error(fmt.Errorf("%v: fetching policy asset: %w", errCtx, err))
		return
	}

	policies, err := io.ReadAll(r)
	if err != nil {
		out.Error(fmt.Errorf("%v: reading policy assets: %w", errCtx, err))
		return
	}

	out.SuccessReadPolicy(string(policies))
}

type Option func(*Interactor)

func WithAuth(a Auth) Option {
	return func(i *Interactor) {
		i.Auth = a
	}
}

func New(r Store, opts ...Option) Interactor {
	i := &Interactor{r, auth.NewVoidAuth()}
	for _, opt := range opts {
		opt(i)
	}
	return *i
}
