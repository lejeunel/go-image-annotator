package query

import (
	"testing"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestSucceedsOnEmpty(t *testing.T) {
	p := NewOrderParserBuilder().Build()
	assert.NoError(t, p.Validate(""))
}

func TestFailsOnUnsupportedField(t *testing.T) {
	p := NewOrderParserBuilder().AddField("gender").Build()
	assert.ErrorIs(t, p.Validate("salary"), e.ErrValidation)
}

func TestSucceedsWithOrderSuffix(t *testing.T) {
	p := NewOrderParserBuilder().AddField("amount").Build()
	assert.NoError(t, p.Validate("amount:asc"))
	assert.NoError(t, p.Validate("amount:desc"))
}

func TestFailsWithUnsupportedSuffix(t *testing.T) {
	p := NewOrderParserBuilder().AddField("amount").Build()
	assert.ErrorIs(t, p.Validate("amount:whatever"), e.ErrValidation)
}

func TestParseOrderingWithEmpty(t *testing.T) {
	p := NewOrderParserBuilder().Build()
	q, err := p.Parse("")
	assert.NoError(t, err)
	assert.Equal(t, im.OrderingArgs{}, q)
}

func TestParseOrderingWithInvalidField(t *testing.T) {
	p := NewOrderParserBuilder().Build()
	_, err := p.Parse("non-existing-field")
	assert.Error(t, err)
}

func TestParseAscendingOrdering(t *testing.T) {
	p := NewOrderParserBuilder().AddField("amount").Build()
	q, err := p.Parse("amount")
	assert.NoError(t, err)
	assert.Equal(t, im.OrderingArgs{{Field: "amount", Order: im.AscOrder}}, q)
}

func TestParseAscendingOrderingWithSuffix(t *testing.T) {
	p := NewOrderParserBuilder().AddField("amount").Build()
	q, err := p.Parse("amount:asc")
	assert.NoError(t, err)
	assert.Equal(t, im.OrderingArgs{{Field: "amount", Order: im.AscOrder}}, q)
}

func TestParseDescendingOrderingWithSuffix(t *testing.T) {
	p := NewOrderParserBuilder().AddField("amount").Build()
	q, err := p.Parse("amount:desc")
	assert.NoError(t, err)
	assert.Equal(t, im.OrderingArgs{{Field: "amount", Order: im.DescOrder}}, q)
}

func TestParseMultipleOrderings(t *testing.T) {
	b := NewOrderParserBuilder()
	b.AddField("amount").Build()
	p := b.AddField("age").Build()
	q, err := p.Parse("amount:desc age:asc")
	assert.NoError(t, err)
	assert.Equal(t, im.OrderingArgs{{Field: "amount", Order: im.DescOrder},
		{Field: "age", Order: im.AscOrder}}, q)
}

func TestParseWithRegexp(t *testing.T) {
	b := NewOrderParserBuilder()
	p := b.AddRegExpField(`^subject\..*$`).Build()
	q, err := p.Parse("subject.age:desc")
	assert.NoError(t, err)
	assert.Equal(t, im.OrderingArgs{{Field: "subject.age", Order: im.DescOrder}}, q)
}

func TestParseWithRenameRule(t *testing.T) {
	b := NewOrderParserBuilder()
	b.AddField(`date-of-birth`)
	b.AddRegExpField(`^meta\..*$`)
	b.AddRenameRule(`\bmeta\.(.*)\b`, `json_extract(m.meta, '$.$1')`)
	p := b.Build()
	q, err := p.Parse("meta.age:desc date-of-birth")
	assert.NoError(t, err)
	assert.Equal(t,
		im.OrderingArgs{
			{Field: "json_extract(m.meta, '$.age')", Order: im.DescOrder},
			{Field: "date-of-birth", Order: im.AscOrder},
		}, q)
}
