package authorizer

import (
	"fmt"
	"slices"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type Policies map[string][]string

func (p Policies) Validate() error {
	invalidMethods := []string{}
	for _, methods := range p {
		for _, method := range methods {
			if !slices.Contains(ValidMethods, method) {
				invalidMethods = append(invalidMethods, method)
			}
		}
	}
	if len(invalidMethods) > 0 {
		return fmt.Errorf(
			"validating methods: found invalid names %v: %w",
			invalidMethods,
			e.ErrValidation,
		)
	}
	return nil
}

var DefaultPolicyFileName = "policies.yaml"

var DefaultPolicies = Policies{
	"viewer":    {},
	"annotator": {"Annotate"},
	"image-contributor": {
		"IngestImage",
		"ImportImage",
		"CreateCollection",
		"CloneCollection",
		"DeleteCollection",
	},
	"admin": {"*"},
}
