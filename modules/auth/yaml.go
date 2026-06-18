package auth

import (
	"fmt"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"gopkg.in/yaml.v2"
	"io"
)

type YamlConfigAuthRule struct {
	Method      string   `yaml:"method"`
	IgnoreGroup bool     `yaml:"ignore_group"`
	Roles       []string `yaml:"roles"`
}

type YamlConfigAuthRules struct {
	Rules []YamlConfigAuthRule `yaml:"rules"`
}

func validateYamlAuthRules(rules []YamlConfigAuthRule) error {
	invalidNames := []string{}
	for _, r := range rules {
		_, ok := validMethods[r.Method]
		if !ok {
			invalidNames = append(invalidNames, r.Method)
		}

	}
	if len(invalidNames) > 0 {
		return fmt.Errorf("checking for validity of method names: got invalid names: %v: %w",
			invalidNames, e.ErrValidation)
	}
	return nil
}

func NewAuthRulesFromYaml(r io.Reader) (*[]AuthRule, error) {
	errCtx := "loading authorization rules from yaml file"
	data, err := io.ReadAll(r)
	if err != nil {
		panic(fmt.Errorf("%v: %w", errCtx, err))
	}
	var yamlAuthRules YamlConfigAuthRules
	if err := yaml.Unmarshal(data, &yamlAuthRules); err != nil {
		return nil, fmt.Errorf("%v: %w: %w", errCtx, err, e.ErrValidation)
	}
	if err := validateYamlAuthRules(yamlAuthRules.Rules); err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
	}

	rules := []AuthRule{}
	for _, r := range yamlAuthRules.Rules {
		rules = append(rules, AuthRule{
			Method:      r.Method,
			IgnoreGroup: r.IgnoreGroup,
			Roles:       r.Roles},
		)
	}

	return &rules, nil
}
