package query

import (
	"fmt"
	"strings"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	s "github.com/lejeunel/go-image-annotator/shared/sql"
	"go.tomakado.io/dumbql"
	"go.tomakado.io/dumbql/schema"
)

type FilterStrParser interface {
	Parse(string) (s.SQLizer, error)
}

type FilterParser struct {
	Schema            schema.Schema
	FieldNameMappings map[string]string
}

type FilterParserOption func(*FilterParser)

func WithFieldNameMapping(source, destination string) FilterParserOption {
	return func(p *FilterParser) {
		p.FieldNameMappings[source] = destination
	}
}

func NewFilterParser(schm schema.Schema, opts ...FilterParserOption) FilterParser {
	p := &FilterParser{Schema: schm, FieldNameMappings: make(map[string]string)}
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

type MappedSQLizer struct {
	s.SQLizer
	Mappings map[string]string
}

func (m MappedSQLizer) ToSql() (string, []any, error) {
	sql, args, err := m.SQLizer.ToSql()
	if err != nil {
		return "", nil, err
	}
	for src, dst := range m.Mappings {
		sql = strings.ReplaceAll(sql, src+" =", dst+" =")
	}
	return sql, args, nil
}

func (v FilterParser) Parse(query string) (s.SQLizer, error) {
	expr, err := dumbql.Parse(query)
	if err != nil {
		return nil, err
	}
	return MappedSQLizer{expr, v.FieldNameMappings}, nil
}
