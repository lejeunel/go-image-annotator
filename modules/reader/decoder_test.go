package reader

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoverEncodedBytes(t *testing.T) {
	decoder := NewBase64ImageDecoder(
		base64.StdEncoding.EncodeToString(testJPGImage),
	)
	r, err := io.ReadAll(decoder)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(r, testJPGImage))
}
