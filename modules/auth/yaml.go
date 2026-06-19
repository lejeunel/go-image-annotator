package auth

import (
	"fmt"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

type YamlConfigAuthRule struct {
	Method      string   `yaml:"method"`
	IgnoreGroup bool     `yaml:"ignore_group"`
	Roles       []string `yaml:"roles"`
	AdminOnly   bool     `yaml:"admin_only"`
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
			Roles:       r.Roles,
			AdminOnly:   r.AdminOnly,
		},
		)
	}

	return &rules, nil
}

func ReadAuthRulesFromPath(path string) (*[]AuthRule, error) {
	voidRules := []AuthRule{}
	if path == "" {
		return &voidRules, nil
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
		return &voidRules, nil
	}
	return rules, nil

}
