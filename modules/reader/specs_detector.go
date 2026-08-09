package reader

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"slices"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func formatToMIME(format string) string {
	switch format {
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

type ImageSpecsDetector struct {
	allowedMIMETypes []string
}

func NewImageSpecsDetector(allowedMIMETypes []string) ImageSpecsDetector {
	return ImageSpecsDetector{allowedMIMETypes}
}

func (d ImageSpecsDetector) Detect(r io.Reader) (*im.Specs, io.Reader, error) {
	// Read a small prefix (DecodeConfig does not need the full file)
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	cfg, format, err := image.DecodeConfig(tee)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding image: %w: %w", err, e.ErrValidation)
	}

	// Reconstruct the full reader:
	// first the consumed bytes, then the remaining original reader
	newReader := io.MultiReader(&buf, r)

	mime := formatToMIME(format)
	if !slices.Contains(d.allowedMIMETypes, mime) {
		return nil, nil, fmt.Errorf("checking whether detected mimetype %v is contained in provided set %v: %w",
			mime, d.allowedMIMETypes, e.ErrValidation)
	}

	return &im.Specs{
		MIMEType: formatToMIME(format),
		Width:    cfg.Width, Height: cfg.Height,
	}, newReader, nil
}
