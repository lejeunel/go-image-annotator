package reader

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"io"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrOnInvalidDataShouldFail(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"jpeg"}, "invalid-data")
	_, err := io.ReadAll(decoder)
	assert.ErrorIs(t, err, e.ErrImageFormat)
}

func TestDecodeJPGImage(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"jpeg"}, base64.StdEncoding.EncodeToString(testJPGImage))
	_, err := io.ReadAll(decoder)
	assert.NoError(t, err)
}

func TestFormatNotAllowedShouldFail(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"png"}, base64.StdEncoding.EncodeToString(testJPGImage))
	_, err := io.ReadAll(decoder)
	assert.ErrorIs(t, err, e.ErrImageFormat)
}

func RecoverEncodedBytes(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"jpeg"}, base64.StdEncoding.EncodeToString(testJPGImage))
	r, err := io.ReadAll(decoder)
	assert.NoError(t, err)
	assert.True(t, bytes.Equal(r, testJPGImage))

}
