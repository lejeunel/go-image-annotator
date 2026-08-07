package authorizer

import (
	"fmt"
	"io"
	"maps"
	"os"
	"strings"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"gopkg.in/yaml.v3"
)

type YamlPolicies struct {
	Version int                 `yaml:"version"`
	Rules   map[string][]string `yaml:"rules"`
}

func MarshalPolicies(policies Policies, w io.Writer) error {
	out := YamlPolicies{Version: 1, Rules: make(map[string][]string)}
	maps.Copy(out.Rules, policies)

	enc := yaml.NewEncoder(w)
	defer enc.Close()

	enc.SetIndent(4)

	return enc.Encode(out)
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
	rules := make(Policies)
	for role, methods := range yamlAuthRules.Rules {
		for _, method := range methods {
			rules[role] = append(rules[role], method)
		}
	}
	if err := rules.Validate(); err != nil {
		return nil, fmt.Errorf("%v: %w", errCtx, err)
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
