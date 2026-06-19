package token_generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// A Base64 encoding processes data in chunks of 3 bytes (24 bits) and converts each chunk into 4 characters (4 × 6 bits).
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
