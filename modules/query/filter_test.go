package query

import (
	"github.com/stretchr/testify/assert"
	"go.tomakado.io/dumbql/schema"
	"testing"
)

func Setup() FilterParser {
	b := schema.NewSchemaBuilder()
	b.AddField("collection", schema.Is[string]())
	return NewFilterParser(b.Build())

}

func TestParse(t *testing.T) {
	p := Setup()
	expr, err := p.Parse("collection:a-collection")
	assert.NoError(t, err)
	sql, _, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "collection = ?", sql)
}

func TestParseWithFieldNameMapping(t *testing.T) {
	p := Setup()
	p.FieldNameMappings = map[string]string{"collection": "collections.name"}
	expr, err := p.Parse("collection:\"a-collection\"")
	assert.NoError(t, err)
	sql, args, err := expr.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "collections.name = ?", sql)
	assert.Equal(t, "a-collection", args[0])
}
