package set

import (
	"context"
	"fmt"
	"io"
	"strings"

	a "github.com/lejeunel/go-image-annotator/modules/authorizer"
)

type Store interface {
	Store(string, io.Reader) error
}

type Auth interface {
	SetPolicies(ctx context.Context) error
	SetAuthRules(rules a.Policies)
}

type Interactor struct {
	Store
	Auth
}

func (i *Interactor) Execute(ctx context.Context, policies string, out OutputPort) {
	errCtx := "setting access policies"
	if err := i.Auth.SetPolicies(ctx); err != nil {
		out.Error(fmt.Errorf("%v: %w", errCtx, err))
		return
	}

	rules, err := a.NewAuthRulesFromYaml(strings.NewReader(policies))
	if err != nil {
		out.Error(fmt.Errorf("%v: reading rules from yaml payload: %w", errCtx, err))
		return
	}

	i.Auth.SetAuthRules(*rules)

	out.SuccessSetPolicy(string(policies))
}

func New(r Store, a Auth) Interactor {
	return Interactor{r, a}
}
