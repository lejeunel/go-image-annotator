package authorizer

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"gopkg.in/yaml.v2"
)

type YamlPolicies struct {
	Version int                 `yaml:"version"`
	Rules   map[string][]string `yaml:"rules"`
}

type Policies map[string][]string

func validateYamlPolicies(cfg YamlPolicies) error {
	invalidNames := []string{}
	for _, methods := range cfg.Rules {
		for _, method := range methods {
			if !slices.Contains(validMethods, method) {
				invalidNames = append(invalidNames, method)
			}
		}
	}
	if len(invalidNames) > 0 {
		return fmt.Errorf("checking for validity of method names: got invalid names: %v: %w",
			invalidNames, e.ErrValidation)
	}
	return nil
}

func NewAuthRulesFromYaml(r io.Reader) (*Policies, error) {
	errCtx := "loading authorization rules from yaml file"
	data, err := io.ReadAll(r)
	if err != nil {
		panic(fmt.Errorf("%v: %w", errCtx, err))
	}
	var yamlAuthRules YamlPolicies
	if err := yaml.Unmarshal(data, &yamlAuthRules); err != nil {
		return nil, fmt.Errorf("%v: %w: %w", errCtx, err, e.ErrValidation)
	}
	if err := validateYamlPolicies(yamlAuthRules); err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	rules := make(Policies)
	for role, methods := range yamlAuthRules.Rules {
		for _, method := range methods {
			if !slices.Contains(validMethods, method) {
				return nil, fmt.Errorf("%v: %w", errCtx, err)

			}
			rules[role] = append(rules[role], method)
		}
	}

	return &rules, nil
}

func ReadAuthRulesFromPath(path string) (*Policies, error) {
	voidPolicies := Policies{}
	if path == "" {
		return &voidPolicies, nil
	}
	errCtx := fmt.Errorf("parsing authentication specifications from file %v", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: file does not exist", errCtx)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: reading file: %w", errCtx, err)
	}
	rules, err := NewAuthRulesFromYaml(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCtx, err)
	}
	if rules == nil {
		return &voidPolicies, nil
	}
	return rules, nil

}
