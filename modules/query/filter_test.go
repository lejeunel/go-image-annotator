package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.tomakado.io/dumbql/schema"
)

func TestParse(t *testing.T) {
	b := NewFilterParserBuilder()
	b.AddField("collection", schema.Is[string]())
	p := b.Build()
	expr, err := p.Parse("collection:\"a-collection\"")
	assert.NoError(t, err)
	sql, args, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "collection = ?", sql)
	assert.Equal(t, "a-collection", args[0])
}

func TestParseWithFieldNameMapping(t *testing.T) {
	b := NewFilterParserBuilder()
	b.AddField("collection", schema.Is[string]())
	b.AddRenameRule(`\bcollection\b`, `collections.name`)
	p := b.Build()

	expr, err := p.Parse("collection:\"a-collection\"")
	assert.NoError(t, err)
	sql, args, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "collections.name = ?", sql)
	assert.Equal(t, "a-collection", args[0])
}

func TestParseWithJSONExtractMapping(t *testing.T) {
	b := NewFilterParserBuilder()
	b.AddRegExpField(`^meta\..*$`, schema.Is[string]())
	b.AddRenameRule(`\bmeta\.(.*)\b`, `json_extract(metadata.meta, '$1')`)
	p := b.Build()

	expr, err := p.Parse("meta.name:\"a-name\"")
	assert.NoError(t, err)
	sql, args, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, `json_extract(metadata.meta, 'name') = ?`, sql)
	assert.Equal(t, "a-name", args[0])
}
func TestDocumentedFilteringField(t *testing.T) {
	b := NewFilterParserBuilder()
	field := "the-field"
	description := "the-description"
	b.AddField(field, schema.Is[string](), WithDescription(description))
	p := b.Build()
	assert.Equal(t, 1, len(p.DescribeFields()))
	assert.Equal(t, field, p.DescribeFields()[0].Name)
	assert.Equal(t, description, p.DescribeFields()[0].Description)
}

func TestDocumentedFilteringRegexpField(t *testing.T) {
	b := NewFilterParserBuilder()
	field := "the-field"
	description := "the-description"
	b.AddRegExpField(field, schema.Is[string](), WithDescription(description))
	p := b.Build()
	assert.Equal(t, 1, len(p.DescribeFields()))
	assert.Equal(t, field, p.DescribeFields()[0].Name)
	assert.Equal(t, description, p.DescribeFields()[0].Description)
}

// func TestExampleFiltering(t *testing.T) {
// 	b := NewOrderParserBuilder()
// 	example := "ingested_at:desc"
// 	b.AddExample(example)
// 	p := b.Build()
// 	assert.Contains(t, p.Examples(), example)
// }
