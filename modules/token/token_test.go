package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTokenOfCorrectLength(t *testing.T) {
	length := 3
	gen := New(length)
	pair, _ := gen.Generate()
	assert.Equal(t, 4, len(pair.Value))
}

func TestVerify(t *testing.T) {
	gen := New(3)
	pair, _ := gen.Generate()
	assert.True(t, gen.Verify(pair.Value, pair.Hash))
}
