package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.tomakado.io/dumbql/query"
	"go.tomakado.io/dumbql/schema"
)

func Setup() FilterParser {
	b := schema.NewSchemaBuilder()
	b.AddField("collection", schema.Is[string]())
	return NewFilterParser(b.Build())
}
func TestParse(t *testing.T) {
	p := Setup()
	expr, err := p.Parse("collection:\"a-collection\"")
	assert.NoError(t, err)
	sql, args, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "collection = ?", sql)
	assert.Equal(t, "a-collection", args[0])
}

func TestParseWithFieldNameMapping(t *testing.T) {
	sb := schema.NewSchemaBuilder()
	sb.AddField("collection", schema.Is[string]())
	rb := query.NewRenamerBuilder()
	rb.Add(`\bcollection\b`, `collections.name`)
	p := NewFilterParser(sb.Build(), WithRenamer(rb.Build()))

	expr, err := p.Parse("collection:\"a-collection\"")
	assert.NoError(t, err)
	sql, args, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "collections.name = ?", sql)
	assert.Equal(t, "a-collection", args[0])
}

func TestParseWithJSONExtractMapping(t *testing.T) {
	sb := schema.NewSchemaBuilder()
	sb.AddRegExpField(`^meta\..*$`, schema.Is[string]())
	rb := query.NewRenamerBuilder()
	rb.Add(`\bmeta\.(.*)\b`, `json_extract(meta_table.meta, '$1')`)
	p := NewFilterParser(sb.Build(), WithRenamer(rb.Build()))

	expr, err := p.Parse("meta.name:\"a-name\"")
	assert.NoError(t, err)
	sql, args, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, `json_extract(metadata.meta, 'name') = ?`, sql)
	assert.Equal(t, "a-name", args[0])
}
