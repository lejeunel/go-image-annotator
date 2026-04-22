package reader

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"errors"
	"io"
	"testing"

	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

func TestErrOnInvalidDataShouldFail(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"jpeg"}, "invalid-data")
	_, err := io.ReadAll(decoder)
	if !errors.Is(err, e.ErrImageFormat) {
		t.Fatalf("expected image format error on invalid base64 data, got %v", err)
	}
}

func TestDecodeJPGImage(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"jpeg"}, base64.StdEncoding.EncodeToString(testJPGImage))
	_, err := io.ReadAll(decoder)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
}

func TestFormatNotAllowedShouldFail(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"png"}, base64.StdEncoding.EncodeToString(testJPGImage))
	_, err := io.ReadAll(decoder)
	if !errors.Is(err, e.ErrImageFormat) {
		t.Fatalf("expected error on invalid format got %v", err)
	}
}

func RecoverEncodedBytes(t *testing.T) {
	decoder := NewBase64ImageDecoder([]string{"jpeg"}, base64.StdEncoding.EncodeToString(testJPGImage))
	r, err := io.ReadAll(decoder)
	if err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	if !bytes.Equal(r, testJPGImage) {
		t.Fatal("did not recover original byte slice")
	}

}
