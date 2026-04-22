package reader

import (
	"bytes"
	"encoding/base64"
	"fmt"
	e "github.com/lejeunel/go-image-annotator-v2/shared/errors"
	i "image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"slices"
)

type Base64ImageDecoder struct {
	AllowedFormats []string
	Base64Data     string
	data           []byte
	offset         int
}

func NewBase64ImageDecoder(allowedFormats []string, base64Data string) *Base64ImageDecoder {
	return &Base64ImageDecoder{AllowedFormats: allowedFormats, Base64Data: base64Data}

}

func (r *Base64ImageDecoder) Read(p []byte) (int, error) {
	errCtx := "decoding base64 data"

	if r.data == nil {
		bytesData, err := base64.StdEncoding.DecodeString(r.Base64Data)
		if err != nil {
			return 0, fmt.Errorf("%v: %v: %w", errCtx, err, e.ErrImageFormat)
		}

		_, format, err := i.Decode(bytes.NewBuffer(bytesData))
		if err != nil {
			return 0, fmt.Errorf("%v: decoding image data: %v: %w", errCtx, err, e.ErrImageFormat)
		}

		if !slices.Contains(r.AllowedFormats, format) {
			return 0, fmt.Errorf("%v: checking for supported format (allowed formats are %v): %v: %w",
				errCtx, r.AllowedFormats, err, e.ErrImageFormat)
		}
		r.data = bytesData

	}
	if r.offset >= len(r.data) {
		return 0, io.EOF
	}

	n := copy(p, r.data[r.offset:])
	r.offset += n

	return n, nil
}
