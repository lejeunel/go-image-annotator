package query

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

type RenameTask struct {
	rex      *regexp.Regexp
	template string
}

type OrderParser struct {
	Fields      []string
	RegExps     []*regexp.Regexp
	RenameTasks []RenameTask
}

type OrderParserBuilder struct {
	fields  []string
	regexps []*regexp.Regexp
	renames []RenameTask
}

func NewOrderParserBuilder() *OrderParserBuilder {
	return &OrderParserBuilder{}
}

func (o *OrderParserBuilder) AddField(field string) *OrderParserBuilder {
	o.fields = append(o.fields, field)
	return o
}
func (o *OrderParserBuilder) AddRegExpField(r string) *OrderParserBuilder {
	re := regexp.MustCompile(r)
	o.regexps = append(o.regexps, re)
	return o
}
func (o *OrderParserBuilder) AddRenameRule(rex string, template string) *OrderParserBuilder {
	re := regexp.MustCompile(rex)
	o.renames = append(o.renames, RenameTask{re, template})
	return o

}

func (o *OrderParserBuilder) Build() OrderParser {
	return OrderParser{Fields: o.fields, RegExps: o.regexps, RenameTasks: o.renames}
}

func (v OrderParser) Validate(q string) error {
	if q == "" {
		return nil
	}
	terms := strings.SplitSeq(q, " ")
	for term := range terms {
		field := v.extractField(term)
		errField := v.validateOnFields(field)
		errRegex := v.validateOnRegExps(field)
		if (errField == nil) || (errRegex == nil) {
			if err := v.validateSuffix(term); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("checking for field %v in known fields: %w", field, e.ErrValidation)

		}
	}
	return nil
}
func (v OrderParser) validateOnFields(field string) error {
	if !slices.Contains(v.Fields, field) {
		return e.ErrValidation
	}
	return nil
}

func (v OrderParser) validateOnRegExps(field string) error {
	for _, re := range v.RegExps {
		if re.MatchString(string(field)) {
			return nil
		}
	}
	return e.ErrValidation
}

func (v OrderParser) Parse(q string) (im.OrderingArgs, error) {
	if q == "" {
		return im.OrderingArgs{}, nil
	}
	if err := v.Validate(q); err != nil {
		return im.OrderingArgs{}, err
	}

	var res im.OrderingArgs
	terms := strings.SplitSeq(q, " ")
	for term := range terms {
		if strings.Contains(term, ":") {
			splits := strings.Split(term, ":")
			field, suffix := splits[0], splits[1]
			field = v.applyRenaming(field)
			r := im.OrderingArg{Field: field}
			if suffix == "desc" {
				r.Order = im.DescOrder
			} else {
				r.Order = im.AscOrder
			}
			res = append(res, r)
		} else {
			res = append(res, im.OrderingArg{Field: v.applyRenaming(term), Order: im.AscOrder})
		}
	}
	return res, nil
}
func (v OrderParser) applyRenaming(field string) string {
	for _, t := range v.RenameTasks {
		field = t.rex.ReplaceAllString(field, t.template)
	}
	return field
}

func (v OrderParser) extractField(term string) string {
	if strings.Contains(term, ":") {
		return strings.Split(term, ":")[0]
	}
	return term
}

func (v OrderParser) validateSuffix(term string) error {
	if strings.Contains(term, ":") {
		suffix := strings.Split(term, ":")[1]
		if !slices.Contains([]string{"asc", "desc"}, suffix) {
			return fmt.Errorf(
				"validating suffix %v to be one of {asc, desc}: %w",
				suffix,
				e.ErrValidation,
			)
		}
	}
	return nil
}
