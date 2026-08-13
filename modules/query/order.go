package query

import (
	"fmt"
	"slices"
	"strings"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type OrderingStrConverter struct {
	Fields []string
}

type OrderingValidatorOption func(*OrderingStrConverter)

func WithOrderingField(field string) OrderingValidatorOption {
	return func(v *OrderingStrConverter) {
		v.Fields = append(v.Fields, field)
	}
}

func NewOrderingConverter(opts ...OrderingValidatorOption) OrderingStrConverter {
	v := &OrderingStrConverter{}
	for _, opt := range opts {
		opt(v)
	}
	return *v
}

func (v OrderingStrConverter) Validate(q string) error {
	if q == "" {
		return nil
	}
	terms := strings.SplitSeq(q, " ")
	for term := range terms {
		field := v.extractField(term)
		if !slices.Contains(v.Fields, field) {
			return fmt.Errorf("checking for field %v in known fields: %w", field, e.ErrValidation)
		}
		if err := v.validateSuffix(term); err != nil {
			return err
		}
	}
	return nil
}
func (v OrderingStrConverter) Parse(q string) (string, error) {
	if q == "" {
		return "", nil
	}
	if err := v.Validate(q); err != nil {
		return "", err
	}

	var res []string
	terms := strings.SplitSeq(q, " ")
	for term := range terms {
		if strings.Contains(term, ":") {
			splits := strings.Split(term, ":")
			field, suffix := splits[0], splits[1]
			if suffix == "desc" {
				res = append(res, field+" DESC")
			} else {
				res = append(res, field)
			}
		} else {
			res = append(res, term)

		}
	}
	return strings.Join(res, ", "), nil

}

func (v OrderingStrConverter) extractField(term string) string {
	if strings.Contains(term, ":") {
		return strings.Split(term, ":")[0]
	}
	return term
}

func (v OrderingStrConverter) validateSuffix(term string) error {
	if strings.Contains(term, ":") {
		suffix := strings.Split(term, ":")[1]
		if !slices.Contains([]string{"asc", "desc"}, suffix) {
			return fmt.Errorf("validating suffix %v to be one of {asc, desc}: %w", suffix, e.ErrValidation)
		}
	}
	return nil
}
