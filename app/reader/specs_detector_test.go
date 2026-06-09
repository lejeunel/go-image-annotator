package reader

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestDetectImageTypeFromNonImageBytesShouldFail(t *testing.T) {
	detector := ImageSpecsDetector{}
	_, _, err := detector.Detect(bytes.NewBuffer([]byte("asdf")))
	assert.ErrorIs(t, err, e.ErrValidation)
}

func TestDetectImageFromJPGBytes(t *testing.T) {
	detector := ImageSpecsDetector{}
	got, _, _ := detector.Detect(bytes.NewBuffer(testJPGImage))
	assert.Equal(t, "image/jpeg", got.MIMEType)
}

func TestDetectImageMIMEType(t *testing.T) {
	detector := ImageSpecsDetector{}
	got, _, _ := detector.Detect(bytes.NewBuffer(testPNGImage))
	assert.Equal(t, "image/png", got.MIMEType)
}

func TestDetectWidth(t *testing.T) {
	detector := ImageSpecsDetector{}
	got, _, _ := detector.Detect(bytes.NewBuffer(testPNGImage))
	assert.Equal(t, 400, got.Width)
}

func TestRecoverBytesAfterDetection(t *testing.T) {
	detector := ImageSpecsDetector{}
	_, reader, _ := detector.Detect(bytes.NewBuffer(testPNGImage))
	r, _ := io.ReadAll(reader)
	assert.True(t, bytes.Equal(r, testPNGImage))
}
