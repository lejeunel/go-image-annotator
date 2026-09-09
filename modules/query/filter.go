package query

import (
	"fmt"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"go.tomakado.io/dumbql"
	"go.tomakado.io/dumbql/query"
	"go.tomakado.io/dumbql/schema"
)

type FilterSQLizer interface {
	ParseToSql(string) (*SQLizer, error)
}

type FilterParserBuilder struct {
	fields        []Field
	descriptions  []string
	schemaBuilder schema.SchemaBuilder
	renameBuilder query.FieldRenamerBuilder
}

func NewFilterParserBuilder() FilterParserBuilder {
	return FilterParserBuilder{schemaBuilder: schema.NewSchemaBuilder(),
		renameBuilder: query.NewRenamerBuilder()}
}

func (b *FilterParserBuilder) AddField(field FieldName, rule RuleFunc, opts ...FieldOption) *FilterParserBuilder {
	new := Field{Name: field}
	for _, opt := range opts {
		opt(&new)
	}
	b.fields = append(b.fields, new)
	b.schemaBuilder.AddField(schema.Field(field), rule)
	return b
}
func (b *FilterParserBuilder) AddRegExpField(field FieldName, rule RuleFunc, opts ...FieldOption) *FilterParserBuilder {
	new := Field{Name: field}
	for _, opt := range opts {
		opt(&new)
	}
	b.fields = append(b.fields, new)
	b.schemaBuilder.AddRegExpField(field, rule)
	return b
}

func (b *FilterParserBuilder) AddRenameRule(rex string, template string) *FilterParserBuilder {
	b.renameBuilder.Add(rex, template)
	return b
}

func (b *FilterParserBuilder) Build() FilterParser {
	return NewFilterParser(b.schemaBuilder.Build(), b.fields, WithRenamer(b.renameBuilder.Build()))
}

type FilterParser struct {
	Schema       schema.Schema
	FieldRenamer *query.FieldRenamer
	Fields       []Field
}

type FilterParserOption func(*FilterParser)

func WithRenamer(renamer query.FieldRenamer) FilterParserOption {
	return func(p *FilterParser) {
		p.FieldRenamer = &renamer
	}
}

func NewFilterParser(schm schema.Schema, fields []Field, opts ...FilterParserOption) FilterParser {
	p := &FilterParser{Schema: schm, Fields: fields}
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
func (v FilterParser) DescribeFilteringFields() []FieldDescription {
	descriptions := []FieldDescription{}
	for _, f := range v.Fields {
		descriptions = append(descriptions, FieldDescription{Name: f.Name, Description: f.Description})
	}
	return descriptions

}
func (v FilterParser) Examples() []string {
	return []string{"first-example", "second-example"}
}
