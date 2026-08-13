package query

import (
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestSucceedsOnEmpty(t *testing.T) {
	v := NewOrderingConverter()
	assert.NoError(t, v.Validate(""))
}

func TestFailsOnUnsupportedField(t *testing.T) {
	v := NewOrderingConverter(WithOrderingField("sex"))
	assert.ErrorIs(t, v.Validate("gender"), e.ErrValidation)
}

func TestSucceedsWithOrderSuffix(t *testing.T) {
	v := NewOrderingConverter(WithOrderingField("amount"))
	assert.NoError(t, v.Validate("amount:asc"))
	assert.NoError(t, v.Validate("amount:desc"))
}

func TestFailsWithUnsupportedSuffix(t *testing.T) {
	v := NewOrderingConverter(WithOrderingField("amount"))
	assert.ErrorIs(t, v.Validate("amount:whatever"), e.ErrValidation)
}

func TestParseOrderingWithEmpty(t *testing.T) {
	v := NewOrderingConverter()
	q, err := v.Parse("")
	assert.NoError(t, err)
	assert.Equal(t, "", q)
}

func TestParseOrderingWithInvalidField(t *testing.T) {
	v := NewOrderingConverter()
	_, err := v.Parse("non-existing-field")
	assert.Error(t, err)
}

func TestParseAscendingOrdering(t *testing.T) {
	v := NewOrderingConverter(WithOrderingField("amount"))
	q, err := v.Parse("amount")
	assert.NoError(t, err)
	assert.Equal(t, "amount", q)
}

func TestParseAscendingOrderingWithSuffix(t *testing.T) {
	v := NewOrderingConverter(WithOrderingField("amount"))
	q, err := v.Parse("amount:asc")
	assert.NoError(t, err)
	assert.Equal(t, "amount", q)
}

func TestParseDescendingOrderingWithSuffix(t *testing.T) {
	v := NewOrderingConverter(WithOrderingField("amount"))
	q, err := v.Parse("amount:desc")
	assert.NoError(t, err)
	assert.Equal(t, "amount DESC", q)
}

func TestParseMultipleOrderings(t *testing.T) {
	v := NewOrderingConverter(WithOrderingField("amount"), WithOrderingField("age"))
	q, err := v.Parse("amount:desc age:asc")
	assert.NoError(t, err)
	assert.Equal(t, "amount DESC, age", q)
}
