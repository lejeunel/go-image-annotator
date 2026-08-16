package query

import (
	"fmt"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"go.tomakado.io/dumbql"
	"go.tomakado.io/dumbql/query"
	"go.tomakado.io/dumbql/schema"
)

type IFilterParser interface {
	Parse(string) (query.Expr, error)
	Validate(query string) error
	ParseToSql(string) (*SQLizer, error)
}

type FilterSQLizer interface {
	ParseToSql(string) (*SQLizer, error)
}

type FilterParser struct {
	Schema       schema.Schema
	FieldRenamer *query.FieldRenamer
}

type FilterParserOption func(*FilterParser)

func WithRenamer(renamer query.FieldRenamer) FilterParserOption {
	return func(p *FilterParser) {
		p.FieldRenamer = &renamer
	}
}

func NewFilterParser(schm schema.Schema, opts ...FilterParserOption) FilterParser {
	p := &FilterParser{Schema: schm}
	for _, opt := range opts {
		opt(p)
	}
	return *p
}

func (v FilterParser) Validate(query string) error {
	expr, err := dumbql.Parse(query)
	if err != nil {
		return fmt.Errorf("validating query %v: %v: %w", query, err, e.ErrValidation)
	}

	_, err = expr.Validate(v.Schema)
	return err
}

func (v FilterParser) Parse(q string) (query.Expr, error) {
	expr, err := dumbql.Parse(q)
	if err != nil {
		return nil, err
	}

	if v.FieldRenamer != nil {
		v.FieldRenamer.Rename(expr)
	}
	return expr, nil
}

func (v FilterParser) ParseToSql(q string) (*SQLizer, error) {
	expr, err := v.Parse(q)
	if err != nil {
		return nil, err
	}
	sql, args, err := expr.ToSql()
	if err != nil {
		return nil, err
	}
	sqlizer := NewSQLizer(sql, args)
	return &sqlizer, nil
}
