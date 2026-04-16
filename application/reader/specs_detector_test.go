package reader

import (
	"bytes"
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"testing"

	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
)

func TestDetectImageTypeFromNonImageBytesShouldFail(t *testing.T) {
	detector := ImageSpecsDetector{}
	_, _, err := detector.Detect(bytes.NewBuffer([]byte("asdf")))
	if !errors.Is(err, e.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestDetectImageFromJPGBytes(t *testing.T) {
	detector := ImageSpecsDetector{}
	got, _, _ := detector.Detect(bytes.NewBuffer(testJPGImage))
	if got.MIMEType != "image/jpeg" {
		t.Fatalf("expected mimetype image/jpg, got %v", got.MIMEType)
	}
}

func TestDetectImageMIMEType(t *testing.T) {
	detector := ImageSpecsDetector{}
	got, _, _ := detector.Detect(bytes.NewBuffer(testPNGImage))
	if got.MIMEType != "image/png" {
		t.Fatalf("expected mimetype image/png, got %v", got.MIMEType)
	}
}

func TestDetectWidth(t *testing.T) {
	detector := ImageSpecsDetector{}
	got, _, _ := detector.Detect(bytes.NewBuffer(testPNGImage))
	want := 400
	if got.Width != want {
		t.Fatalf("expected width of %v, got %v", want, got.Width)
	}
}

func TestRecoverBytesAfterDetection(t *testing.T) {
	detector := ImageSpecsDetector{}
	_, reader, _ := detector.Detect(bytes.NewBuffer(testPNGImage))
	r, _ := io.ReadAll(reader)
	if !bytes.Equal(r, testPNGImage) {
		t.Fatal("expected to retrieve original bytes after detection")
	}

}
